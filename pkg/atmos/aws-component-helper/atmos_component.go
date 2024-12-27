package aws_component_helper

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

type AtmosComponent struct {
	RandomIdentifier string
	ComponentName    string
	StackName        string
}

func NewAtmosComponent(component string, stack string) *AtmosComponent {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	return &AtmosComponent{
		RandomIdentifier: randomId,
		ComponentName:    component,
		StackName:        stack,
	}
}

func (ac *AtmosComponent) getAtmosOptions(t *testing.T, options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone()
	}

	currentTFDataDir := ".terraform"
	if value, ok := options.EnvVars["TF_DATA_DIR"]; ok {
		currentTFDataDir = value
	}
	stack := strings.Replace(ac.StackName, "/", "-", -1)
	name := strings.Replace(ac.ComponentName, "/", "-", -1)
	envvars := map[string]string{
		// We need to split the TF_DATA_DIR for parallel suites mode
		"TF_DATA_DIR":             filepath.Join(currentTFDataDir, fmt.Sprintf("component-%s", ac.RandomIdentifier)),
		"TEST_WORKSPACE_TEMPLATE": fmt.Sprintf("%s-%s-%s", stack, name, ac.RandomIdentifier),
	}

	err := mergo.Merge(&result.EnvVars, envvars)
	require.NoError(t, err)

	// Merge in any additional vars passed in
	err = mergo.Merge(&result.Vars, vars)
	require.NoError(t, err)

	result.Component = ac.ComponentName
	result.Stack = ac.StackName

	atmosOptions := atmos.WithDefaultRetryableErrors(t, result)
	return atmosOptions
}
