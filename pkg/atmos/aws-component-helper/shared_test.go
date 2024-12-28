package aws_component_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchFilter(t *testing.T) {
	t.Parallel()
	t.Run("wrong", func(t *testing.T) {
		matchRegexp := ""
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance")
		assert.Error(t, err)
		assert.False(t, result)
	})

	t.Run("default", func(t *testing.T) {
		matchRegexp := ""
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("suite", func(t *testing.T) {
		matchRegexp := "default"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("all tests", func(t *testing.T) {
		matchRegexp := "default"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default/test1")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("one test tests", func(t *testing.T) {
		matchRegexp := "default/test1"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default/test1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("multilevel subtests", func(t *testing.T) {
		matchRegexp := "default/test1"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default/test1/subtest1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test1/subtest2")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest1")
		assert.NoError(t, err)
		assert.False(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("filter longer then test", func(t *testing.T) {
		matchRegexp := "default/test1/subtest1"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default/test1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest1")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("regexp", func(t *testing.T) {
		matchRegexp := "default/.*/subtest1"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default/test1/subtest1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test1/subtest2")
		assert.NoError(t, err)
		assert.False(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

}
