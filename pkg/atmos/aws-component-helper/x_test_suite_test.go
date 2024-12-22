package aws_component_helper

import (
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentTestSuiteMinimum(t *testing.T) {
	componentTestSuite := NewXTestSuite(&atmos.Options{})
	test := componentTestSuite.GetOrCreateTest("test")
	test.SetSubject("vpc", "default-test")

	assert.Equal(t, test.Subject.Component, "vpc")
	assert.Equal(t, componentTestSuite.tests["test"].Subject.Component, "vpc")
	assert.Equal(t, componentTestSuite.tests["test"].Subject.Stack, "default-test")
}
