package aws_component_helper

import (
	"dario.cat/mergo"
	"flag"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	skipDeploySuiteDependencies  = flag.Bool("cth.skip-deploy-suite-deps", false, "skip deploy suite deps")
	skipDestroySuiteDependencies = flag.Bool("cth.skip-destroy-suite-deps", false, "skip destroy suite deps")
)

type XTestSuite struct {
	RandomIdentifier string
	setup            []*AtmosComponent
	tests            map[string]*ComponentTest
	atmosOptions     *atmos.Options
}

func NewXTestSuite() *XTestSuite {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	return &XTestSuite{
		RandomIdentifier: randomId,
		setup:            make([]*AtmosComponent, 0),
		tests:            make(map[string]*ComponentTest),
	}
}

func (ts *XTestSuite) getAtmosOptions(t *testing.T, options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone()
	}

	err := mergo.Merge(&result.Vars, vars)
	require.NoError(t, err)

	suiteVars := map[string]interface{}{
		"attributes": []string{ts.RandomIdentifier},
		"default_tags": map[string]string{
			"CreatedByAtmosTestSuite": ts.RandomIdentifier,
		},
	}

	err = mergo.Merge(&result.Vars, suiteVars)
	require.NoError(t, err)
	return result
}

func (ts *XTestSuite) AddSetup(component string, stack string) {
	item := NewAtmosComponent(component, stack)
	ts.setup = append(ts.setup, item)
}

func (ts *XTestSuite) GetOrCreateTest(name string) *ComponentTest {
	if _, ok := ts.tests[name]; !ok {
		ts.tests[name] = NewComponentTest()
	}
	return ts.tests[name]
}

func (ts *XTestSuite) Tests() map[string]*ComponentTest {
	return ts.tests
}

func (ts *XTestSuite) Run(t *testing.T, options *atmos.Options) {
	suiteOptions := ts.getAtmosOptions(t, options, map[string]interface{}{})
	for _, component := range ts.setup {
		componentOptions := component.getAtmosOptions(t, suiteOptions, map[string]interface{}{})
		if !*skipDeploySuiteDependencies {
			atmosApply(t, componentOptions)
		}
		if !*skipDeploySuiteDependencies && !*skipDestroySuiteDependencies {
			defer atmosDestroy(t, componentOptions)
		}
	}

	if *runParallel {
		fmt.Println("Run tests in parallel mode")
		t.Parallel()
	}
	for name, item := range ts.tests {
		t.Run(name, func(t *testing.T) {
			item.Run(t, suiteOptions)
		})
	}
}
