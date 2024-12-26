package aws_component_helper

import (
	"dario.cat/mergo"
	"flag"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
	"testing"
)

var (
	skipVerifyEnabledFlag = flag.Bool("cth.skip-verify-enabled-flag", true, "skip verify enabled flag")
)

type ComponentTest struct {
	RandomIdentifier string
	setup            []*AtmosComponent
	Subject          *AtmosComponent
	assert           map[string]func(t *testing.T, ct *ComponentTest)
}

func NewComponentTest() *ComponentTest {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	return &ComponentTest{
		RandomIdentifier: randomId,
		setup:            make([]*AtmosComponent, 0),
		Subject:          nil,
		assert:           map[string]func(t *testing.T, ct *ComponentTest){},
	}
}

func (ct *ComponentTest) verifyEnabledFlag(t *testing.T, component *AtmosComponent, options *atmos.Options) *atmos.Options {
	testOptions := ct.getAtmosOptions(t, options, map[string]interface{}{})
	vars := map[string]interface{}{
		"enabled": false,
	}

	componentOptions := component.getAtmosOptions(t, testOptions, vars)

	exitCode, err := atmosPlanExitCodeE(t, componentOptions)
	require.NoError(t, err)

	if exitCode != 0 {
		require.Fail(t, "running atmos terraform plan with enabled flag set to false resulted in resource changes")
	}

	return options
}

func (ct *ComponentTest) getAtmosOptions(t *testing.T, options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone()
	}

	currentTFDataDir := ".terraform"
	if value, ok := options.EnvVars["TF_DATA_DIR"]; ok {
		currentTFDataDir = value
	}

	envvars := map[string]string{
		// We need to split the TF_DATA_DIR for parallel suites mode
		"TF_DATA_DIR": filepath.Join(currentTFDataDir, fmt.Sprintf("test-%s", ct.RandomIdentifier)),
	}

	err := mergo.Merge(&result.EnvVars, envvars)
	require.NoError(t, err)

	mergedVars := map[string]interface{}{
		"default_tags": map[string]string{
			"CreatedByAtmosTestSuiteTest": ct.RandomIdentifier,
		},
	}

	// Merge in any additional vars passed in
	err = mergo.Merge(&result.Vars, mergedVars)
	require.NoError(t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(t, err)

	return result
}

func (ct *ComponentTest) AddSetup(component string, stack string) {
	item := NewAtmosComponent(component, stack)
	ct.setup = append(ct.setup, item)
}

func (ct *ComponentTest) SetSubject(component string, stack string) {
	ct.Subject = NewAtmosComponent(component, stack)
}

func (ct *ComponentTest) AddSAssert(name string, callback func(t *testing.T, ct *ComponentTest)) {
	ct.assert[name] = callback
}
