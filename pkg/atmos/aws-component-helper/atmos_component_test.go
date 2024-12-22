package aws_component_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtmosComponent(t *testing.T) {
	// Create a temporary atmos component
	component := NewAtmosComponent("vpc", "default-test")

	assert.Equal(t, component.Component, "vpc")
	assert.Equal(t, component.Stack, "default-test")
}
