package examples_helper

import (
	"context"
	"github.com/charmbracelet/log"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"testing"
)

type LocalStackConfiguration struct {
	Services []string
	Image    string

	HostPort            string // Set by the localstack container run
	LocalStackContainer testcontainers.Container

	bSetStsEndpointToLocalStack bool // Set AWS StS endpoint to localstack when created
}

func NewLocalStackConfiguration() *LocalStackConfiguration {
	return &LocalStackConfiguration{
		Services: []string{"s3", "iam", "lambda", "dynamodb", "sts", "account", "ec2"},
		Image:    "localstack/localstack:1.4.0",

		bSetStsEndpointToLocalStack: true,
	}
}

func (s *TestSuite) SetupLocalStackContainer(t *testing.T, config *c.Config) {
	//if s.Config.SkipSetupLocalStack {
	//	return
	//}

	ctx := context.Background()
	// Create a slice with the required ports
	ports := []int{4566}

	// Append ports from 4510 to 4559
	for i := 4510; i <= 4559; i++ {
		ports = append(ports, i)
	}

	LocalStackContainer, err := localstack.Run(
		ctx, "localstack/localstack:4.2.0",
		testcontainers.WithEnv(map[string]string{
			"SERVICES":              "s3, iam, lambda, dynamodb, sts, account, ec2, kms",
			"DEBUG":                 "1",
			"DOCKER_HOST":           "unix:///var/run/docker.sock",
			"AWS_ACCESS_KEY_ID":     "test",
			"AWS_SECRET_ACCESS_KEY": "test",
		}),
		testcontainers.WithHostPortAccess(ports...),
	)
	s.SetupConfiguration.LocalStackConfiguration.LocalStackContainer = LocalStackContainer

	portMap, err := LocalStackContainer.Ports(ctx)
	if err != nil {
		log.WithPrefix(t.Name()).Fatal("Failed to get ports from container", "error", err)
	}

	hostPort := portMap[nat.Port("4566/tcp")][0].HostPort //  [{HostIP:0.0.0.0 HostPort:56614}]
	s.SetupConfiguration.LocalStackConfiguration.HostPort = hostPort
	t.Setenv("LOCALSTACK_PORT", hostPort)
	if s.SetupConfiguration.LocalStackConfiguration.bSetStsEndpointToLocalStack {
		// Used by awsutils and is required for dependencies
		t.Setenv("AWS_ENDPOINT_URL_STS", "https://localhost.localstack.cloud:"+hostPort)
	}
	assert.NoError(t, err, "failed to start localstack container")

	s.logPhaseStatus("setup/localstack container", "completed")
}
