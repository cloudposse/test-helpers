package aws_component_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentTestMinimum(t *testing.T) {
	componentTest := NewComponentTest()
	componentTest.SetSubject("vpc", "default-test")

	assert.Equal(t, componentTest.Subject.Component, "vpc")
	assert.Equal(t, componentTest.Subject.Stack, "default-test")
}
