package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAwsComponentTestOptions(t *testing.T) {
	// Test default values
	options := NewAwsComponentTestOptions("us-west-2", "test-component", "test-stack")

	assert.Equal(t, "us-west-2", options.AwsRegion)
	assert.Equal(t, "test-component", options.ComponentName)
	assert.Equal(t, "test-stack", options.StackName)
	assert.Equal(t, "./fixtures", options.FixturesPath)
	assert.False(t, options.SkipAwsNuke)
	assert.Empty(t, options.StackDependencies)

	// Test with custom options
	customOptions := NewAwsComponentTestOptions(
		"us-east-1",
		"custom-component",
		"custom-stack",
		WithFixturesPath("/custom/path"),
		WithSkipAwsNuke(true),
		WithDependencies([]Stack{{Component: "dep1", StackName: "stack1"}}),
	)

	assert.Equal(t, "us-east-1", customOptions.AwsRegion)
	assert.Equal(t, "custom-component", customOptions.ComponentName)
	assert.Equal(t, "custom-stack", customOptions.StackName)
	assert.Equal(t, "/custom/path", customOptions.FixturesPath)
	assert.True(t, customOptions.SkipAwsNuke)
	assert.Len(t, customOptions.StackDependencies, 1)
	assert.Equal(t, "dep1", customOptions.StackDependencies[0].Component)
	assert.Equal(t, "stack1", customOptions.StackDependencies[0].StackName)
}

func TestWithFixturesPath(t *testing.T) {
	options := &AwsComponentTestOptions{}
	WithFixturesPath("/test/path")(options)
	assert.Equal(t, "/test/path", options.FixturesPath)
}

func TestWithSkipAwsNuke(t *testing.T) {
	options := &AwsComponentTestOptions{}
	WithSkipAwsNuke(true)(options)
	assert.True(t, options.SkipAwsNuke)
}

func TestWithDependencies(t *testing.T) {
	options := &AwsComponentTestOptions{}
	dependencies := []Stack{
		{Component: "dep1", StackName: "stack1"},
		{Component: "dep2", StackName: "stack2"},
	}
	WithDependencies(dependencies)(options)
	assert.Equal(t, dependencies, options.StackDependencies)
}
