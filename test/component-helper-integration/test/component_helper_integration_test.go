package test

import (
	"fmt"
	"testing"

	atmos "github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	"github.com/stretchr/testify/assert"
)

type ExampleTestSuite struct {
	helper.TestSuite
}

func (s *ExampleTestSuite) TestEnabledFlag() {
	s.VerifyEnabledFlag("component2", "test-use2-sandbox", nil)
}

func (s *ExampleTestSuite) TestComponentTestHelperExample1() {
	const component = "component2"
	const stack = "test-use2-sandbox"

	defer s.DestroyAtmosComponent(s.T(), component, stack, nil)

	options, _ := s.DeployAtmosComponent(s.T(), component, stack, nil)

	revision := atmos.Output(s.T(), options, "revision")

	expectedRevision := fmt.Sprintf("%s-%s", s.Config.RandomIdentifier, "2")
	assert.Equal(s.T(), expectedRevision, revision)
}

func (s *ExampleTestSuite) TestComponentTestHelperExample2() {
	assert.Equal(s.T(), 2, 2)
}

func (s *ExampleTestSuite) TestComponentTestHelperExample3() {
	assert.Equal(s.T(), 3, 3)
}

func TestRunExampleSuite(t *testing.T) {
	suite := new(ExampleTestSuite)

	suite.AddDependency(t, "component1", "test-use2-sandbox", nil)
	suite.AddDependency(t, "component3", "test-use2-sandbox", nil)
	helper.Run(t, suite)
}
