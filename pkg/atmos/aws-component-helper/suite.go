package aws_component_helper

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	skipTests = flag.Bool("skip-tests", false, "skip tests")
	skipSetup = flag.Bool("skip-setup", false, "skip setup")
)

type teadDown struct {
	component *AtmosComponent
	callback  *func(t *testing.T, atm *Atmos)
}

type Suite struct {
	t                *testing.T
	randomIdentifier string
	name             string
	dependencies     []*AtmosComponent
	teardown         []*teadDown
	options          *atmos.Options
	state            *State
}

func NewSuite(t *testing.T, name string, fixture *Fixture) *Suite {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	require.NotContains(t, name, "/")
	suiteState, err := fixture.State.Fork(name)
	require.NoError(t, err)
	suite := &Suite{
		t:                t,
		name:             name,
		randomIdentifier: randomId,
		dependencies:     []*AtmosComponent{},
		state:            suiteState,
		teardown:         []*teadDown{},
		options:          fixture.getAtmosOptions(&atmos.Options{}, map[string]interface{}{}),
	}
	return suite
}

func (ts *Suite) AddDependency(componentName string, stackName string) {
	component := NewAtmosComponent(componentName, stackName, nil)
	ts.dependencies = append(ts.dependencies, component)
	ts.teardown = append(ts.teardown, &teadDown{component: component, callback: nil})
	if *skipSetup {
		fmt.Printf("Skip suite %s setup dependency component: %s stack: %s\n", ts.name, componentName, stackName)
		return
	}
	ts.getAtmos(ts.state).Deploy(component)
}

func (ts *Suite) getAtmos(state *State) *Atmos {
	return NewAtmos(ts.t, state, ts.getAtmosOptions(map[string]interface{}{
		"attributes": []string{ts.randomIdentifier},
	}))
}

func (ts *Suite) getTestAtmos(state *State) *Atmos {
	return NewAtmos(ts.t, state, ts.getAtmosOptions(map[string]interface{}{}))
}

func (ts *Suite) runTeardown() {
	if *skipTeardown {
		fmt.Printf("Skip teardown suite %s\n", ts.name)
		return
	}
	atm := ts.getAtmos(ts.state)
	var f *teadDown
	for i := len(ts.teardown) - 1; i >= 0; i-- {
		f = ts.teardown[i]
		if f.callback != nil {
			(*f.callback)(ts.t, atm)
		}
		if f.component != nil {
			atm.Destroy(f.component)
		}
	}
	err := ts.state.Teardown()
	assert.NoError(ts.t, err)
}

func (ts *Suite) Setup(t *testing.T, f func(t *testing.T, atm *Atmos)) {
	if *skipSetup {
		fmt.Printf("Skip suite %s setup callback\n", ts.name)
		return
	}
	atm := ts.getAtmos(ts.state)
	f(t, atm)
}

func (ts *Suite) TearDown(t *testing.T, f func(t *testing.T, atm *Atmos)) {
	ts.teardown = append(ts.teardown, &teadDown{component: nil, callback: &f})
}

func (ts *Suite) Test(t *testing.T, name string, f func(t *testing.T, atm *Atmos)) {
	if *skipTests {
		fmt.Printf("Skip test %s/%s\n", ts.name, name)
		return
	}

	testState, err := ts.state.Fork(name)
	require.NoError(t, err)
	defer testState.Teardown()

	atm := ts.getTestAtmos(testState)
	testRunName := fmt.Sprintf("%s/%s", t.Name(), name)
	if ok, err := matchFilter(testRunName); ok {
		t.Run(name, func(t *testing.T) {
			f(t, atm)
		})
	} else {
		require.NoError(t, err)
	}
}

func (ts *Suite) getAtmosOptions(vars map[string]interface{}) *atmos.Options {
	result, err := ts.options.Clone()
	require.NoError(ts.t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(ts.t, err)

	return result
}

func (ts *Suite) GetRandomIdentifier() string {
	return ts.randomIdentifier
}
