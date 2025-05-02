package examples_helper

import (
	"bufio"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"io"
	"os"
	"strings"
	"testing"
)

type LocalStackConfiguration struct {
	Services []string
	Image    string

	HostPort            string // Set by the localstack container run
	LocalStackContainer testcontainers.Container

	bUpdateAWSEndpointsToLocalStack bool // Set AWS StS endpoint to localstack when created

	UseDockerComposeInstance bool // Set to true if using docker compose instance
}

func NewLocalStackConfiguration() *LocalStackConfiguration {
	return &LocalStackConfiguration{
		Services: []string{"s3", "iam", "lambda", "dynamodb", "sts", "account", "ec2"},
		Image:    "localstack/localstack:4.2.0",

		bUpdateAWSEndpointsToLocalStack: true,
	}
}
func streamContainerLogs(ctx context.Context, container testcontainers.Container, logger *log.Logger) {
	logs, err := container.Logs(ctx)
	if err != nil {
		logger.Error("Failed to retrieve container logs", "error", err)
		return
	}
	defer logs.Close()

	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		logger.Info(scanner.Text(), "source", "LocalStack")
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		logger.Error("Error reading container logs", "error", err)
	}
}

func (s *TestSuite) SetupLocalStackContainer(t *testing.T, config *c.Config) {
	if s.SetupConfiguration.LocalStackConfiguration.UseDockerComposeInstance {
		s.SetupConfiguration.LocalStackConfiguration.HostPort = "4566"
		t.Setenv("LOCALSTACK_PORT", "4566")
		s.UpdateAwsEnvVarsToLocalStack(t)
		return
	}
	ctx := context.Background()

	ports := createLocalStackPortArray()

	if s.Config.SkipTearDownLocalStack {
		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	}

	//logger := getLocalStackLogger(t)
	log.WithPrefix(t.Name()).Info("Starting localstack container", "image", s.SetupConfiguration.LocalStackConfiguration.Image)

	LocalStackContainer, err := localstack.Run(ctx,
		s.SetupConfiguration.LocalStackConfiguration.Image,
		//testcontainers.WithLogger(logger),

		testcontainers.WithEnv(map[string]string{
			"SERVICES":              s.getLocalStackServices(),
			"DEBUG":                 "1",
			"DOCKER_HOST":           "unix:///var/run/docker.sock",
			"AWS_ACCESS_KEY_ID":     "test",
			"AWS_SECRET_ACCESS_KEY": "test",
			"LOCALSTACK_AUTH_TOKEN": os.Getenv("LOCALSTACK_AUTH_TOKEN"),
		}),
		testcontainers.WithHostPortAccess(ports...),
	)
	// Capture and log container logs
	//go streamContainerLogs(ctx, LocalStackContainer, logger)
	s.SetupConfiguration.LocalStackConfiguration.LocalStackContainer = LocalStackContainer

	portMap, err := LocalStackContainer.Ports(ctx)
	if err != nil {
		log.WithPrefix(t.Name()).Error("Failed to get ports from container", "error", err)
	}

	hostPort := portMap[nat.Port("4566/tcp")][0].HostPort //  [{HostIP:0.0.0.0 HostPort:56614}]
	s.SetupConfiguration.LocalStackConfiguration.HostPort = hostPort
	t.Setenv("LOCALSTACK_PORT", hostPort)
	if s.SetupConfiguration.LocalStackConfiguration.bUpdateAWSEndpointsToLocalStack {
		// Used by awsutils and is required for dependencies
		s.UpdateAwsEnvVarsToLocalStack(s.T())
	}
	assert.NoError(t, err, "failed to start localstack container")

	s.logPhaseStatus("setup/localstack container", "completed")
}

func (s *TestSuite) getLocalStackServices() string {
	return strings.Join(s.SetupConfiguration.LocalStackConfiguration.Services, ", ")
}

func createLocalStackPortArray() []int {
	ports := []int{4566}

	// Append ports from 4510 to 4559
	for i := 4510; i <= 4559; i++ {
		ports = append(ports, i)
	}
	return ports
}

func getLocalStackLogger(t *testing.T) *log.Logger {
	logger := log.Default().WithPrefix(t.Name()).WithPrefix("localstack")
	logger.SetLevel(log.DebugLevel)
	testcontainers.Logger = logger
	return logger
}

func (s *TestSuite) UpdateAwsEnvVarsToLocalStack(t *testing.T) {
	hostport := s.SetupConfiguration.LocalStackConfiguration.HostPort
	t.Setenv("AWS_REGION", "us-east-1")
	//https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/custom-service-endpoints#available-endpoint-customizations
	// AWS Backend Variables
	// https://developer.hashicorp.com/terraform/language/v1.5.x/settings/backends/s3#configuration
	localhostConfig := "http://localhost:" + hostport
	localstackCloudConfig := "https://localhost.localstack.cloud:" + hostport
	localstackS3Endpoint := "http://s3.localhost.localstack.cloud:" + hostport

	t.Setenv("AWS_S3_ENDPOINT", localstackS3Endpoint)
	t.Setenv("AWS_DYNAMODB_ENDPOINT", localstackCloudConfig)
	t.Setenv("AWS_STS_ENDPOINT", localhostConfig)
	t.Setenv("AWS_ENDPOINT_URL", localhostConfig)

	t.Setenv("AWS_ENDPOINT_URL_STS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ACCESSANALYZER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ACCOUNT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ACM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ACM_PCA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_AMP", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_AMPLIFY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_API_GATEWAY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APIGATEWAYV2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPLICATION_AUTO_SCALING", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPCONFIG", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPFABRIC", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPFLOW", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPINTEGRATIONS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPLICATION_INSIGHTS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPLICATION_SIGNALS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APP_MESH", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPRUNNER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPSTREAM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_APPSYNC", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ATHENA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_AUDITMANAGER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_AUTO_SCALING", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_AUTO_SCALING_PLANS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_BACKUP", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_BATCH", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_BCM_DATA_EXPORTS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_BEDROCK", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_BEDROCK_AGENT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_BILLING", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_BUDGETS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_COST_EXPLORER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CHATBOT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CHIME", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CHIME_SDK_MEDIA_PIPELINES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CHIME_SDK_VOICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLEANROOMS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUD9", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDCONTROL", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDFORMATION", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDFRONT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDFRONT_KEYVALUESTORE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDHSM_V2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDSEARCH", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDTRAIL", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDWATCH", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODEARTIFACT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODEBUILD", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODECATALYST", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODECOMMIT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODECONNECTIONS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODEGURUPROFILER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODEGURU_REVIEWER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODEPIPELINE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODESTAR_CONNECTIONS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODESTAR_NOTIFICATIONS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_COGNITO_IDENTITY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_COGNITO_IDENTITY_PROVIDER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_COMPREHEND", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_COMPUTE_OPTIMIZER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CONFIG_SERVICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CONNECT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CONNECTCASES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CONTROLTOWER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_COST_OPTIMIZATION_HUB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_COST_AND_USAGE_REPORT_SERVICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CUSTOMER_PROFILES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DATABREW", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DATAEXCHANGE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DATA_PIPELINE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DATASYNC", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DATAZONE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DAX", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CODEDEPLOY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DETECTIVE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DEVICE_FARM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DEVOPS_GURU", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DIRECT_CONNECT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DLM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DATABASE_MIGRATION_SERVICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DOCDB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DOCDB_ELASTIC", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DRS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DIRECTORY_SERVICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DSQL", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_DYNAMODB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EC2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ECR", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ECR_PUBLIC", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ECS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EFS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EKS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ELASTICACHE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ELASTIC_BEANSTALK", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ELASTICSEARCH_SERVICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ELASTIC_TRANSCODER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ELASTIC_LOAD_BALANCING", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ELASTIC_LOAD_BALANCING_V2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EMR", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EMR_CONTAINERS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EMR_SERVERLESS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EVENTBRIDGE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_EVIDENTLY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_FINSPACE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_FIREHOSE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_FIS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_FMS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_FSX", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GAMELIFT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GLACIER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GLOBAL_ACCELERATOR", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GLUE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GRAFANA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GREENGRASS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GROUNDSTATION", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_GUARDDUTY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_HEALTHLAKE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IAM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IDENTITYSTORE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IMAGEBUILDER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_INSPECTOR", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_INSPECTOR2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_INTERNETMONITOR", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_INVOICING", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IOT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IOTANALYTICS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IOT_EVENTS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IVS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_IVSCHAT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KAFKA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KAFKACONNECT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KENDRA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KEYSPACES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KINESIS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KINESIS_ANALYTICS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KINESIS_ANALYTICS_V2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KINESIS_VIDEO", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_KMS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LAKEFORMATION", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LAMBDA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LAUNCH_WIZARD", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LEX_MODEL_BUILDING_SERVICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LEX_MODELS_V2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LICENSE_MANAGER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LIGHTSAIL", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LOCATION", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_CLOUDWATCH_LOGS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_LOOKOUTMETRICS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_M2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MACIE2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEDIACONNECT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEDIACONVERT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEDIALIVE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEDIAPACKAGE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEDIAPACKAGEV2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEDIAPACKAGE_VOD", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEDIASTORE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MEMORYDB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MGN", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MQ", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_MWAA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_NEPTUNE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_NEPTUNE_GRAPH", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_NETWORK_FIREWALL", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_NETWORKMANAGER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_NETWORKMONITOR", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_OAM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_OPENSEARCH", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_OPENSEARCHSERVERLESS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_OPSWORKS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ORGANIZATIONS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_OSIS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_OUTPOSTS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_PAYMENTCRYPTOGRAPHY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_PCA_CONNECTOR_AD", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_PCS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_PINPOINT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_PINPOINT_SMS_VOICE_V2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_PIPES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_POLLY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_PRICING", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_QBUSINESS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_QLDB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_QUICKSIGHT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RAM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RBIN", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RDS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_REDSHIFT", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_REDSHIFT_DATA", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_REDSHIFT_SERVERLESS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_REKOGNITION", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RESILIENCEHUB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RESOURCE_EXPLORER_2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RESOURCE_GROUPS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RESOURCE_GROUPS_TAGGING_API", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ROLESANYWHERE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ROUTE_53", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ROUTE_53_DOMAINS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ROUTE_53_PROFILES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ROUTE53_RECOVERY_CONTROL_CONFIG", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ROUTE53_RECOVERY_READINESS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_ROUTE53RESOLVER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_RUM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_S3", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_S3_CONTROL", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_S3OUTPOSTS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_S3TABLES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SAGEMAKER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SCHEDULER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SCHEMAS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SECRETS_MANAGER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SECURITYHUB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SECURITYLAKE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SERVERLESSAPPLICATIONREPOSITORY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SERVICE_CATALOG", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SERVICE_CATALOG_APPREGISTRY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SERVICEDISCOVERY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SERVICE_QUOTAS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SESV2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SFN", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SHIELD", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SIGNER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SIMPLEDB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SNS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SQS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SSM", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SSM_CONTACTS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SSM_INCIDENTS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SSM_QUICKSETUP", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SSM_SAP", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SSO", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SSO_ADMIN", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_STORAGE_GATEWAY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_STS", localhostConfig) // LOCALHOST
	t.Setenv("AWS_ENDPOINT_URL_SWF", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_SYNTHETICS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_TAXSETTINGS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_TIMESTREAM_INFLUXDB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_TIMESTREAM_QUERY", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_TIMESTREAM_WRITE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_TRANSCRIBE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_TRANSFER", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_VERIFIEDPERMISSIONS", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_VPC_LATTICE", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_WAF", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_WAF_REGIONAL", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_WAFV2", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_WELLARCHITECTED", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_WORKLINK", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_WORKSPACES", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_WORKSPACES_WEB", localstackCloudConfig)
	t.Setenv("AWS_ENDPOINT_URL_XRAY", localstackCloudConfig)
}

func (s *TestSuite) NewLocalstackS3Client() *s3.Client {
	hostport := s.SetupConfiguration.LocalStackConfiguration.HostPort
	// Hardcode or env var your LocalStack S3 endpoint
	endpoint := "http://localhost:" + hostport // default LocalStack edge port
	region := "us-east-1"

	// Manually construct config for LocalStack
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		config.WithBaseEndpoint(endpoint),
	)
	if err != nil {
		log.Errorf("failed to load config: %v", err)
		return nil
	}

	return s3.NewFromConfig(cfg)
}

func (s *TestSuite) ShutDownExistingLocalStackContainer(t *testing.T) {
	s.logPhaseStatus("teardown/localstack container", "started")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Errorf("Unabel to create docker client, please make sure that docker is installed\n%s", err.Error())
		t.Fail()
		return
	}
	list, err := cli.ContainerList(context.Background(), container.ListOptions{})
	for _, c := range list {
		if strings.Contains(c.Image, "localstack") || strings.Contains(c.Image, "testcontainers") {
			log.WithPrefix(t.Name()).Info("Stopping localstack container", "container", c.ID)
			cli.ContainerStop(context.Background(), c.ID, container.StopOptions{})
			cli.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{})
		}
	}
	s.logPhaseStatus("teardown/localstack container", "completed")
}
