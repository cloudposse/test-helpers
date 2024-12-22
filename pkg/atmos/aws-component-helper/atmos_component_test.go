package aws_component_helper

import (
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtmosComponent(t *testing.T) {
	// Create a temporary atmos component
	component := NewAtmosComponent("vpc", "default-test", &atmos.Options{})

	assert.Equal(t, component.Component, "vpc")
	assert.Equal(t, component.Stack, "default-test")
}
