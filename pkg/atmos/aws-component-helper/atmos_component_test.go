package aws_component_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAtmosComponent tests the creation and properties of an AtmosComponent
func TestAtmosComponent(t *testing.T) {
	// Create a temporary AtmosComponent with name "vpc" and stack "default-test"
	// Passing 'nil' for vars as no specific variables are required for this test
	component := NewAtmosComponent("vpc", "default-test", nil)

	// Assert that the ComponentName is correctly set to "vpc"
	assert.Equal(t, component.ComponentName, "vpc")

	// Assert that the StackName is correctly set to "default-test"
	assert.Equal(t, component.StackName, "default-test")
}
