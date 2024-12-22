package aws_component_helper

import (
	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/stretchr/testify/require"
	"testing"
)

type AtmosComponent struct {
	Component    string
	Stack        string
	atmosOptions *atmos.Options
}

func NewAtmosComponent(component string, stack string, options *atmos.Options) *AtmosComponent {
	return &AtmosComponent{
		Component:    component,
		Stack:        stack,
		atmosOptions: options,
	}
}

func (ac *AtmosComponent) getAtmosOptions(t *testing.T, vars map[string]interface{}) *atmos.Options {
	result, _ := ac.atmosOptions.Clone()

	// Merge in any additional vars passed in
	err := mergo.Merge(&result.Vars, vars)
	require.NoError(t, err)

	result.Component = ac.Component
	result.Stack = ac.Component

	atmosOptions := atmos.WithDefaultRetryableErrors(t, result)
	return atmosOptions
}
