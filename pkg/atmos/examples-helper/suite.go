package examples_helper

import (
	"os/exec"
	"testing"

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
	TempContentsCmd         *exec.Cmd
	AtmosBaseDir            string // Base directory for atmos relative to the temp dir
	LocalStackConfiguration *LocalStackConfiguration
	VendorAllComponents     bool
	PullBeforeDeploy        bool
}

func NewSetupConfiguration() SetupConfiguration {
	return SetupConfiguration{
		TempContentsCmd:         nil,
		AtmosBaseDir:            "",
		LocalStackConfiguration: NewLocalStackConfiguration(),
		VendorAllComponents:     false,
		PullBeforeDeploy:        true,
	}
}

type TestSuite struct {
	Config       *c.Config
	Dependencies []*dependency.Dependency
	suite.Suite
	SetupConfiguration SetupConfiguration
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
func (s *TestSuite) AddCustomDependency(t *testing.T, d *dependency.Dependency) {
	s.Dependencies = append(s.Dependencies, d)
}
func (s *TestSuite) AddDependency(t *testing.T, componentName string, stackName string, additionalVars *map[string]interface{}, vendor bool, targets []string, addRandomAttribute bool, args ...string) {
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

	s.SetupConfiguration.LocalStackConfiguration = NewLocalStackConfiguration()
}

// Setup runs the setup phase of the test suite.
func (s *TestSuite) SetupSuite() {
	t := s.T()
	if t == nil {
		panic("SetupSuite called with nil *testing.T, call s.SetT(t) first")
	}

	s.InitConfig()
	config := s.Config

	if s.Config.SkipSetupTestSuite {
		s.logPhaseStatus("setup", "skipped")
	} else {
		s.logPhaseStatus("setup", "started")
	}

	s.BootstrapTempDir(t, config)

	s.CreateTempContents(t, config)
	s.SetupLocalStackContainer(t, config)
	if s.SetupConfiguration.VendorAllComponents {
		s.VendorAllComponents(t, config)
	} else {
		s.PullDependencies(t, config)
	}
	s.DeployDependencies(t, config)

	s.logPhaseStatus("setup", "completed")
}

func (s *TestSuite) TearDownSuite() {
	t := s.T()
	if !s.Config.SkipDestroyDependencies {
		s.DestroyDependencies(t, s.Config)
	}

	if s.Config.SkipTeardownTestSuite {
		s.logPhaseStatus("teardown", "skipped")
		return
	}

	s.DestroyTempDir(t, s.Config)
	s.DestroyConfigFile(t, s.Config)
	s.DestroyLocalStackContainer(t, s.Config)
}
