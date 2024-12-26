package aws_component_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentTestMinimum(t *testing.T) {
	componentTest := NewComponentTest()
	componentTest.SetSubject("vpc", "default-test")

	assert.Equal(t, componentTest.Subject.ComponentName, "vpc")
	assert.Equal(t, componentTest.Subject.StackName, "default-test")
}
