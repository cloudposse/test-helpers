package examples_helper

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/charmbracelet/log"
	gwaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"fmt"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
	"github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/dependency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SetupConfiguration struct {
	TempContentsCmd           *exec.Cmd
	AtmosBaseDir              string // Base directory for atmos relative to the temp dir
	LocalStackConfiguration   *LocalStackConfiguration
	VendorAllComponents       bool
	PullBeforeDeploy          bool
	DeployTfStateBackendStack string // Deploy the tfstate backend, stackName or empty string to skip
}

func NewSetupConfiguration() *SetupConfiguration {
	return &SetupConfiguration{
		TempContentsCmd:           nil,
		AtmosBaseDir:              "",
		LocalStackConfiguration:   NewLocalStackConfiguration(),
		VendorAllComponents:       false,
		PullBeforeDeploy:          true,
		DeployTfStateBackendStack: "core-use1-root",
	}
}

type TestSuite struct {
	Config       *c.Config
	Dependencies []*dependency.Dependency
	suite.Suite
	SetupConfiguration *SetupConfiguration
	SuperUserAccessKey string
	SuperUserSecretKey string
}

type TestingSuite interface {
	GetConfig(t *testing.T) *c.Config
	suite.TestingSuite
}

func NewTestSuite() *TestSuite {
	tsuite := new(TestSuite)
	tsuite.Dependencies = make([]*dependency.Dependency, 0)

	return tsuite
}

func Run(t *testing.T, s TestingSuite) {
	suite.Run(t, s)
}

func (s *TestSuite) GetConfig(t *testing.T) *c.Config {
	assert.NotNil(t, s.Config)
	return s.Config
}
func (s *TestSuite) AddDependency(t *testing.T, d *dependency.Dependency) {
	s.Dependencies = append(s.Dependencies, d)
}
func (s *TestSuite) AddWorkflowDependency(t *testing.T, workflowName string, workflowFile string) {
	s.Dependencies = append(s.Dependencies, &dependency.Dependency{
		WorkflowName: workflowName,
		WorkflowFile: workflowFile,
	})
}

func (s *TestSuite) AddComponentDependency(t *testing.T, componentName string, stackName string, additionalVars *map[string]interface{}, vendor bool, targets []string, addRandomAttribute bool, args ...string) {
	s.Dependencies = append(s.Dependencies, &dependency.Dependency{
		AdditionalVars:     additionalVars,
		ComponentName:      componentName,
		StackName:          stackName,
		Args:               args,
		Vendor:             vendor,
		Targets:            targets,
		AddRandomAttribute: addRandomAttribute,
	})
}

func (s *TestSuite) AddVendorOnlyDependency(t *testing.T, componentName string) {
	s.Dependencies = append(s.Dependencies, &dependency.Dependency{
		ComponentName: componentName,
		VendorOnly:    true,
	})
}

func (s *TestSuite) AddFunctionDependency(t *testing.T, fn func() error) {
	s.Dependencies = append(s.Dependencies, &dependency.Dependency{
		Function: fn,
	})
}

func (s *TestSuite) GetAtmosOptions(componentName string, stackName string, additionalVars *map[string]interface{}) *atmos.Options {
	mergedVars := s.getMergedVars(s.T(), additionalVars)
	return GetAtmosOptions(s.T(), s.Config, componentName, stackName, &mergedVars)
}

func (s *TestSuite) getMergedVars(t *testing.T, additionalVars *map[string]interface{}) map[string]interface{} {

	mergedVars := map[string]interface{}{
		"attributes": []string{s.Config.RandomIdentifier},
	}

	// Merge in any additional vars passed in
	if additionalVars != nil {
		err := mergo.Merge(&mergedVars, additionalVars)
		require.NoError(t, err)
	}

	return mergedVars
}

func (s *TestSuite) DeployAtmosComponentWithOptions(t *testing.T, options *atmos.Options, componentName string, stackName string, additionalVars *map[string]interface{}) (*atmos.Options, string) {
	phaseName := fmt.Sprintf("deploy/atmos component/%s/%s", stackName, componentName)
	if s.SetupConfiguration.PullBeforeDeploy {
		s.logPhaseStatus(phaseName, "started")
		s.pullComponent(t, s.Config, componentName)
		s.logPhaseStatus(phaseName, "completed")
	}

	if s.Config.SkipDeployComponent {
		s.logPhaseStatus(phaseName, "skipped")
		return nil, ""
	}

	s.logPhaseStatus(phaseName, "started")

	mergedVars := s.getMergedVars(t, additionalVars)

	atmosOptions := getAtmosOptionsFromSetupConfiguration(t, s.Config, s.SetupConfiguration, componentName, stackName, &mergedVars, nil)
	atmosOptions.MergeOptions(options)
	output, err := atmos.ApplyE(t, atmosOptions)
	if err != nil {
		s.logPhaseStatus(phaseName, "failed")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")

	return atmosOptions, output
}
func (s *TestSuite) DeployAtmosComponent(t *testing.T, componentName string, stackName string, additionalVars *map[string]interface{}) (*atmos.Options, string) {
	phaseName := fmt.Sprintf("deploy/atmos component/%s/%s", stackName, componentName)
	if s.SetupConfiguration.PullBeforeDeploy {
		s.logPhaseStatus(phaseName, "started")
		s.pullComponent(t, s.Config, componentName)
		s.logPhaseStatus(phaseName, "completed")
	}

	if s.Config.SkipDeployComponent {
		s.logPhaseStatus(phaseName, "skipped")
		return nil, ""
	}

	s.logPhaseStatus(phaseName, "started")

	mergedVars := s.getMergedVars(t, additionalVars)
	atmosOptions := getAtmosOptionsFromSetupConfiguration(t, s.Config, s.SetupConfiguration, componentName, stackName, &mergedVars, nil)

	output, err := atmos.ApplyE(t, atmosOptions)
	if err != nil {
		s.logPhaseStatus(phaseName, "failed")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")

	return atmosOptions, output
}

func (s *TestSuite) DestroyAtmosComponent(t *testing.T, componentName string, stackName string, additionalVars *map[string]interface{}) {
	phaseName := fmt.Sprintf("destroy/atmos component/%s/%s", stackName, componentName)

	if s.Config.SkipDestroyComponent {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	mergedVars := s.getMergedVars(t, additionalVars)
	atmosOptions := getAtmosOptionsFromSetupConfiguration(t, s.Config, s.SetupConfiguration, componentName, stackName, &mergedVars, nil)

	_, err := atmos.DestroyE(t, atmosOptions)
	require.NoError(t, err)

	s.logPhaseStatus(phaseName, "completed")
}

func (s *TestSuite) InitConfig() {
	t := s.T()

	if s.Config == nil {
		config := c.InitConfig(t)
		s.Config = config
	}
	s.SetupConfiguration = NewSetupConfiguration()
	s.SetupConfiguration.LocalStackConfiguration = NewLocalStackConfiguration()
}

// Setup runs the setup phase of the test suite.
func (s *TestSuite) SetupSuite() {
	t := s.T()
	if t == nil {
		panic("SetupSuite called with nil *testing.T, call s.SetT(t) first")
	}

	config := s.Config

	if s.Config.SkipSetupTestSuite {
		s.logPhaseStatus("setup", "skipped")
	} else {
		s.logPhaseStatus("setup", "started")
	}

	s.BootstrapTempDir(t, config)
	s.CreateTempContents(t, config)
	if _, err := os.Stat(config.FixturesDir); err == nil {
		s.logPhaseStatus("fixtures", "started")
		s.copyDirectoryRecursively(t, config.FixturesDir, filepath.Join(config.TempDir, s.SetupConfiguration.AtmosBaseDir))
		s.logPhaseStatus("fixtures", "completed")
	}
	s.SetupLocalStackContainer(t, config)
	if s.SetupConfiguration.VendorAllComponents {
		s.VendorAllComponents(t, config)
	} else {
		s.PullDependencies(t, config)
	}

	if s.SetupConfiguration.DeployTfStateBackendStack != "" {
		s.InitTerraformState(t, s.SetupConfiguration.DeployTfStateBackendStack)
	}

	s.DeployDependencies(t, config)

	s.logPhaseStatus("setup", "completed")
}

func (s *TestSuite) TearDownSuite() {

	t := s.T()
	if !s.Config.SkipTearDownLocalStack {
		defer s.DestroyLocalStackContainer(t, s.Config)
	}
	if s.Config.SkipTeardownTestSuite {
		s.logPhaseStatus("teardown", "skipped")
		return
	}

	defer s.DestroyTempDir(t, s.Config)
	defer s.DestroyConfigFile(t, s.Config)

	if !s.Config.SkipDestroyDependencies {
		s.DestroyDependencies(t, s.Config)
	}
}

func (s *TestSuite) RunAtmosWorkflow(t *testing.T, WorkflowName string, WorkflowFile string) {

	phaseName := fmt.Sprintf("run atmos workflow [%s] file: [%s]", WorkflowName, WorkflowFile)
	s.logPhaseStatus(phaseName, "started")

}

func (s *TestSuite) InitTerraformState(t *testing.T, stack string) {
	phaseName := "init tfstate-backend"
	s.logPhaseStatus(phaseName, "started")

	options := getAtmosOptionsFromSetupConfiguration(t, s.Config, s.SetupConfiguration, "tfstate-backend", stack, nil, nil)
	// Vendor
	s.pullComponent(t, s.Config, options.Component)

	// Deploy the tfstate backend
	_, err := atmos.RunAtmosCommandE(t, options, "terraform", "apply", options.Component, "-var=access_roles_enabled=false", "--stack", options.Stack, "--auto-generate-backend-file=false", "-input=false", "-auto-approve")
	if err != nil {
		s.logPhaseStatus(phaseName, "failed")
		require.NoError(t, err)
	}

	retry.DoWithRetryableErrorsE(t, "Waiting for tfstate bucket", options.RetryableAtmosErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {

		out, err := shell.RunCommandAndGetOutputE(t, shell.Command{
			Command: "aws",
			Args:    []string{"s3", "ls"},
		})

		if err != nil {
			return out, err
		}
		out, err = shell.RunCommandAndGetOutputE(t, shell.Command{
			Command: "aws",
			Args:    []string{"sts", "get-caller-identity"},
		})

		if err != nil {
			return out, err
		}
		return out, nil
	})
	err = s.CreateSuperUser(t)

	s.AssumeSuperUser()
	shell.RunCommandAndGetOutputE(t, shell.Command{
		Command: "aws",
		Args:    []string{"sts", "get-caller-identity"},
	})

	time.Sleep(8 * time.Second)
	shell.RunCommandAndGetOutputE(t, shell.Command{
		Command: "aws",
		Args:    []string{"s3", "ls"},
	})
	atmos.RunAtmosCommandE(t, options, "terraform", "init", options.Component, "-s", stack, "--", "-force-copy")

	_, err = atmos.RunAtmosCommandE(t, options, "terraform", "apply", options.Component, "-var=access_roles_enabled=false", "--stack", options.Stack, "--skip-init", "-input=false", "-auto-approve")
	if err != nil {
		s.logPhaseStatus(phaseName, "failed")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")
}

func (s *TestSuite) CreateSuperUser(t *testing.T) error {
	ctx := context.Background()
	SuperAdminUsername := "SuperAdmin"
	iamClient := gwaws.NewIamClient(t, "us-east-1")
	iamClient.CreateUser(ctx, &iam.CreateUserInput{
		UserName: aws.String(SuperAdminUsername),
	})
	_, err := iamClient.AttachUserPolicy(ctx, &iam.AttachUserPolicyInput{
		UserName:  aws.String(SuperAdminUsername),
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
	})
	resp, err := iamClient.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(SuperAdminUsername),
	})
	if err != nil {
		log.WithPrefix(t.Name()).Fatalf("failed to create access key: %v", err)
	}
	log.WithPrefix(t.Name()).Print("Access key created")
	s.SuperUserAccessKey = *resp.AccessKey.AccessKeyId
	s.SuperUserSecretKey = *resp.AccessKey.SecretAccessKey
	return err
}

func (s *TestSuite) AssumeSuperUser() {
	s.T().Setenv("AWS_ACCESS_KEY_ID", s.SuperUserAccessKey)
	s.T().Setenv("AWS_SECRET_ACCESS_KEY", s.SuperUserSecretKey)
}

func (s *TestSuite) AssumeRootAccount() {
	s.T().Setenv("AWS_ACCESS_KEY_ID", "test")
	s.T().Setenv("AWS_SECRET_ACCESS_KEY", "test")
}
