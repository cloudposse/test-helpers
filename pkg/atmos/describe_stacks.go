package atmos

import (
	"encoding/json"
	"github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/stretchr/testify/require"
)

type DescribeStacksTerraformComponent struct {
	AtmosComponent         string                                    `json:"atmos_component"`
	AtmosStack             string                                    `json:"atmos_stack"`
	AtmosStackFile         string                                    `json:"atmos_stack_file"`
	Backend                interface{}                               `json:"backend"`
	BackendType            string                                    `json:"backend_type"`
	Command                string                                    `json:"command"`
	Component              string                                    `json:"component"`
	Env                    interface{}                               `json:"env"`
	Inheritance            interface{}                               `json:"inheritance"`
	Metadata               interface{}                               `json:"metadata"`
	Overrides              interface{}                               `json:"overrides"`
	RemoteStateBackend     interface{}                               `json:"remote_state_backend"`
	RemoteStateBackendType string                                    `json:"remote_state_backend_type"`
	Settings               *DescribeStacksTerraformComponentSettings `json:"settings"`
	Vars                   interface{}                               `json:"vars"`
	Workspace              string                                    `json:"workspace"`
}

type DescribeStacksTerraformComponentSettings struct {
	Test *DescribeStacksTerraformComponentSettingsTest `json:"test, omitempty"`
}

type DescribeStacksTerraformComponentSettingsTest struct {
	Suite string `json:"suite"`
	Test  string `json:"test"`
	Stage string `json:"stage"`
}

type DescribeStacksTerraform struct {
	Terraform map[string]DescribeStacksTerraformComponent `json:"terraform"`
}

type DescribeStacksComponent struct {
	Components DescribeStacksTerraform `json:"components"`
}

type DescribeStacksOutput struct {
	Stacks map[string]DescribeStacksComponent
}

func DescribeStacks(t testing.TestingT, options *Options) DescribeStacksOutput {
	out, err := DescribeStacksE(t, options)
	require.NoError(t, err)
	return out
}

func DescribeStacksE(t testing.TestingT, options *Options) (DescribeStacksOutput, error) {
	result := DescribeStacksOutput{}
	output, err := RunAtmosCommandAndGetStdoutE(t, options, FormatArgs(options, "describe", "stacks", "--component-types=terraform", "--format=json")...)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(output), &result.Stacks)
	return result, err
}
