package aws_component_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtmosComponent(t *testing.T) {
	// Create a temporary atmos component
	component := NewAtmosComponent("vpc", "default-test", nil)

	assert.Equal(t, component.ComponentName, "vpc")
	assert.Equal(t, component.StackName, "default-test")
}
