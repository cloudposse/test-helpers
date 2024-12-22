package aws_component_helper

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type ComponentTest struct {
	RandomIdentifier  string
	setup             []*AtmosComponent
	Subject           *AtmosComponent
	assert            []*AtmosComponent
	VerifyEnabledFlag bool
	AtmosOptions      *atmos.Options
}

func NewComponentTest(options *atmos.Options) *ComponentTest {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	return &ComponentTest{
		RandomIdentifier:  randomId,
		setup:             make([]*AtmosComponent, 0),
		Subject:           nil,
		assert:            make([]*AtmosComponent, 0),
		VerifyEnabledFlag: true,
		AtmosOptions:      options,
	}
}

func (ct *ComponentTest) Run(t *testing.T) {
	for _, component := range ct.setup {
		componentOptions := component.getAtmosOptions(t, map[string]interface{}{})

		atmosApply(t, componentOptions)
		defer atmosDestroy(t, componentOptions)
	}

	if ct.VerifyEnabledFlag {
		fmt.Println("VerifyEnabledFlag")
		ct.verifyEnabledFlag(t, ct.Subject)
	} else {
		fmt.Println("Skipping VerifyEnabledFlag")
	}

	subjectOptions := ct.Subject.getAtmosOptions(t, map[string]interface{}{})
	atmosApply(t, subjectOptions)
	defer atmosDestroy(t, subjectOptions)

	for _, component := range ct.assert {
		componentOptions := component.getAtmosOptions(t, map[string]interface{}{})

		atmosApply(t, componentOptions)
		defer atmosDestroy(t, componentOptions)
	}
}

func (ct *ComponentTest) verifyEnabledFlag(t *testing.T, component *AtmosComponent) *atmos.Options {
	vars := map[string]interface{}{
		"enabled": false,
	}
	options := component.getAtmosOptions(t, vars)

	exitCode, err := atmosPlanExitCodeE(t, options)
	assert.NoError(t, err)

	if exitCode != 0 {
		assert.Fail(t, "running atmos terraform plan with enabled flag set to false resulted in resource changes")
	}

	return options
}

func (ct *ComponentTest) getAtmosOptions() *atmos.Options {
	result, _ := ct.AtmosOptions.Clone()
	mergedVars := map[string]interface{}{
		"default_tags": map[string]string{
			"CreatedByAtmosTestSuiteTest": ct.RandomIdentifier,
		},
	}
	// Merge in any additional vars passed in
	err := mergo.Merge(&result.Vars, mergedVars)
	if err != nil {
		return nil
	}
	//require.NoError(t, err)

	return result
}

func (ct *ComponentTest) AddSetup(component string, stack string) {
	item := NewAtmosComponent(component, stack, ct.getAtmosOptions())
	ct.setup = append(ct.setup, item)
}

func (ct *ComponentTest) SetSubject(component string, stack string) {
	ct.Subject = NewAtmosComponent(component, stack, ct.getAtmosOptions())
}

func (ct *ComponentTest) AddSAssert(component string, stack string) {
	item := NewAtmosComponent(component, stack, ct.getAtmosOptions())
	ct.assert = append(ct.assert, item)
}
