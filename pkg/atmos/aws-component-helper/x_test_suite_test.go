package aws_component_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentTestSuiteMinimum(t *testing.T) {
	componentTestSuite := NewXTestSuite()
	test := componentTestSuite.GetOrCreateTest("test")
	test.SetSubject("vpc", "default-test")

	assert.Equal(t, test.Subject.ComponentName, "vpc")
	assert.Equal(t, componentTestSuite.tests["test"].Subject.ComponentName, "vpc")
	assert.Equal(t, componentTestSuite.tests["test"].Subject.StackName, "default-test")
}
