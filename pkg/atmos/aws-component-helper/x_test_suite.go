package aws_component_helper

import (
	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"strings"
	"testing"
)

type XTestSuite struct {
	RandomIdentifier string
	setup            []*AtmosComponent
	tests            map[string]*ComponentTest
	atmosOptions     *atmos.Options
}

func NewXTestSuite(options *atmos.Options) *XTestSuite {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	return &XTestSuite{
		RandomIdentifier: randomId,
		setup:            make([]*AtmosComponent, 0),
		tests:            make(map[string]*ComponentTest),
		atmosOptions:     options,
	}
}

func (ts *XTestSuite) getAtmosOptions() *atmos.Options {
	result, _ := ts.atmosOptions.Clone()

	vars := map[string]interface{}{
		"attributes": []string{ts.RandomIdentifier},
		"default_tags": map[string]string{
			"CreatedByAtmosTestSuite": ts.RandomIdentifier,
		},
	}

	err := mergo.Merge(&result.Vars, vars)
	if err != nil {
		return nil
	}
	return result
}

func (ts *XTestSuite) AddSetup(component string, stack string) {
	item := NewAtmosComponent(component, stack, ts.getAtmosOptions())
	ts.setup = append(ts.setup, item)
}

func (ts *XTestSuite) GetOrCreateTest(name string) *ComponentTest {
	if _, ok := ts.tests[name]; !ok {
		ts.tests[name] = NewComponentTest(ts.getAtmosOptions())
	}
	return ts.tests[name]
}

func (ts *XTestSuite) Tests() map[string]*ComponentTest {
	return ts.tests
}

func (ts *XTestSuite) Run(t *testing.T) {
	for _, component := range ts.setup {
		componentOptions := component.getAtmosOptions(t, map[string]interface{}{})
		atmosApply(t, componentOptions)
		defer atmosDestroy(t, componentOptions)
	}

	//t.Parallel()
	for name, item := range ts.tests {
		t.Run(name, func(t *testing.T) {
			item.Run(t)
		})
	}
}
