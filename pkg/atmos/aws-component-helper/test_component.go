package aws_component_helper

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
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

		atmosApply(t, componentOptions)
		defer atmosDestroy(t, componentOptions)
	}

	if !*skipVerifyEnabledFlag {
		fmt.Println("VerifyEnabledFlag")
		ct.verifyEnabledFlag(t, ct.Subject, options)
	} else {
		fmt.Println("Skipping VerifyEnabledFlag")
	}

	subjectOptions := ct.Subject.getAtmosOptions(t, testOptions, map[string]interface{}{})
	atmosApply(t, subjectOptions)
	defer atmosDestroy(t, subjectOptions)

	for _, component := range ct.assert {
		componentOptions := component.getAtmosOptions(t, testOptions, map[string]interface{}{})

		atmosApply(t, componentOptions)
		defer atmosDestroy(t, componentOptions)
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
