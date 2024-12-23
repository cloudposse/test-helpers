package aws_component_helper

import (
	"dario.cat/mergo"
	"flag"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	skipVerifyEnabledFlag = flag.Bool("cth.skip-verify-enabled-flag", true, "skip verify enabled flag")

	skipDeployTestDependencies  = flag.Bool("cth.skip-deploy-test-deps", false, "skip deploy test deps")
	skipDestroyTestDependencies = flag.Bool("cth.skip-destroy-test-deps-teardown", false, "skip destroy test deps")

	skipDeployComponentUnderTest  = flag.Bool("cth.skip-deploy-cut", false, "skip deploy component under test")
	skipDestroyComponentUnderTest = flag.Bool("cth.skip-destroy-cut", false, "skip destroy component under test")

	skipDeployAsserts  = flag.Bool("cth.skip-deploy-asserts", false, "skip deploy component under test")
	skipDestroyAsserts = flag.Bool("cth.skip-destroy-asserts", false, "skip destroy component under test")
)

type ComponentTest struct {
	RandomIdentifier string
	setup            []*AtmosComponent
	Subject          *AtmosComponent
	assert           []*AtmosComponent
}

func NewComponentTest() *ComponentTest {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	return &ComponentTest{
		RandomIdentifier: randomId,
		setup:            make([]*AtmosComponent, 0),
		Subject:          nil,
		assert:           make([]*AtmosComponent, 0),
	}
}

func (ct *ComponentTest) Run(t *testing.T, options *atmos.Options) {
	testOptions := ct.getAtmosOptions(t, options, map[string]interface{}{})
	for _, component := range ct.setup {
		componentOptions := component.getAtmosOptions(t, testOptions, map[string]interface{}{})
		if !*skipDeployTestDependencies {
			atmosApply(t, componentOptions)
		}
		if !*skipDeployTestDependencies && !*skipDestroyTestDependencies {
			defer atmosDestroy(t, componentOptions)
		}
	}

	if !*skipVerifyEnabledFlag {
		fmt.Println("VerifyEnabledFlag")
		ct.verifyEnabledFlag(t, ct.Subject, options)
	} else {
		fmt.Println("Skipping VerifyEnabledFlag")
	}

	subjectOptions := ct.Subject.getAtmosOptions(t, testOptions, map[string]interface{}{})
	if !*skipDeployComponentUnderTest {
		atmosApply(t, subjectOptions)
	}
	if !*skipDeployComponentUnderTest && !*skipDestroyComponentUnderTest {
		defer atmosDestroy(t, subjectOptions)
	}

	for _, component := range ct.assert {
		componentOptions := component.getAtmosOptions(t, testOptions, map[string]interface{}{})

		if !*skipDeployAsserts {
			atmosApply(t, componentOptions)
		}
		if !*skipDeployAsserts && !*skipDestroyAsserts {
			defer atmosDestroy(t, componentOptions)
		}
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

	mergedVars := map[string]interface{}{
		"default_tags": map[string]string{
			"CreatedByAtmosTestSuiteTest": ct.RandomIdentifier,
		},
	}

	// Merge in any additional vars passed in
	err := mergo.Merge(&result.Vars, mergedVars)
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

func (ct *ComponentTest) AddSAssert(component string, stack string) {
	item := NewAtmosComponent(component, stack)
	ct.assert = append(ct.assert, item)
}
