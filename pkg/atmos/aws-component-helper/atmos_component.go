package aws_component_helper

import (
	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/stretchr/testify/require"
	"testing"
)

type AtmosComponent struct {
	Component string
	Stack     string
}

func NewAtmosComponent(component string, stack string) *AtmosComponent {
	return &AtmosComponent{
		Component: component,
		Stack:     stack,
	}
}

func (ac *AtmosComponent) getAtmosOptions(t *testing.T, options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone()
	}

	// Merge in any additional vars passed in
	err := mergo.Merge(&result.Vars, vars)
	require.NoError(t, err)

	result.Component = ac.Component
	result.Stack = ac.Stack

	atmosOptions := atmos.WithDefaultRetryableErrors(t, result)
	return atmosOptions
}
