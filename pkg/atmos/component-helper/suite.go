package component_helper

import (
	"testing"

	"fmt"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/component-helper/config"
	"github.com/cloudposse/test-helpers/pkg/atmos/component-helper/dependency"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	Config       *c.Config
	Dependencies []*dependency.Dependency
	suite.Suite
}

func NewTestSuite() *TestSuite {
	tsuite := new(TestSuite)
	tsuite.Dependencies = make([]*dependency.Dependency, 0)

	return tsuite
}

func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

func (s *TestSuite) AddDependency(t *testing.T, componentName string, stackName string, additionalVars *map[string]interface{}) {
	s.Dependencies = append(s.Dependencies, &dependency.Dependency{
		AdditionalVars: additionalVars,
		ComponentName:  componentName,
		StackName:      stackName,
	})
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

	if s.Config.SkipDeployComponent {
		s.logPhaseStatus(phaseName, "skipped")
		return nil, ""
	}

	s.logPhaseStatus(phaseName, "started")

	mergedVars := s.getMergedVars(t, additionalVars)
	atmosOptions := getAtmosOptions(t, s.Config, componentName, stackName, &mergedVars)

	output, err := atmos.ApplyE(t, atmosOptions)
	if err != nil {
		s.logPhaseStatus(phaseName, "failed")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")

	return atmosOptions, output
}

func (s *TestSuite) DestroyAtmosComponent(t *testing.T, componentName string, stackName string, additionalVars *map[string]interface{}) {
	mergedVars := s.getMergedVars(t, additionalVars)
	atmosOptions := getAtmosOptions(t, s.Config, componentName, stackName, &mergedVars)

	_, err := atmos.DestroyE(t, atmosOptions)
	require.NoError(t, err)
}

// Setup runs the setup phase of the test suite.
func (s *TestSuite) SetupSuite() {
	t := s.T()

	config := c.InitConfig(t)
	s.Config = config

	if s.Config.SkipSetupTestSuite {
		s.logPhaseStatus("setup", "skipped")
	} else {
		s.logPhaseStatus("setup", "started")
	}

	s.BootstrapTempDir(t, config)
	s.CopyComponentToTempDir(t, config)
	s.VendorDependencies(t, config)
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
}
