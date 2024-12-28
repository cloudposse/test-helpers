package aws_component_helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

type teadDown struct {
	component *AtmosComponent
	callback  *func(t *testing.T, atm *Atmos)
}

type Suite struct {
	t                *testing.T
	randomIdentifier string
	name             string
	stateDir         string
	globalStateDir   string
	dependencies     []*AtmosComponent
	teardown         []*teadDown
	options          *atmos.Options
}

func NewSuite(t *testing.T, name string, fixture *Fixture) *Suite {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	stateNamespace := fmt.Sprintf("suite-%s", name)
	suite := &Suite{
		t:                t,
		name:             name,
		randomIdentifier: randomId,
		dependencies:     []*AtmosComponent{},
		stateDir:         filepath.Join(fixture.GlobalStateDir(), stateNamespace),
		globalStateDir:   fixture.GlobalStateDir(),
		teardown:         []*teadDown{},
		options:          fixture.getAtmosOptions(&atmos.Options{}, map[string]interface{}{}),
	}

	if fixture.StateDir() != "" {
		copyDirectoryRecursively(fixture.StateDir(), suite.stateDir)
	} else {
		err := createDir(fixture.GlobalStateDir(), stateNamespace)
		require.NoError(t, err)
	}

	return suite
}

func (ts *Suite) AddDependency(componentName string, stackName string) {
	component := NewAtmosComponent(componentName, stackName)
	ts.dependencies = append(ts.dependencies, component)
	ts.teardown = append(ts.teardown, &teadDown{component: component, callback: nil})
	ts.getAtmos(ts.t).Deploy(ts.t, component)
}

func (ts *Suite) getAtmos(t *testing.T) *Atmos {
	return &Atmos{}
}

func (ts *Suite) getTestAtmos(t *testing.T) *Atmos {
	return &Atmos{}
}

func (ts *Suite) runTeardown(t *testing.T) {
	atm := ts.getAtmos(t)
	var f *teadDown
	for i := len(ts.teardown) - 1; i >= 0; i-- {
		f = ts.teardown[i]
		if f.callback != nil {
			(*f.callback)(t, atm)
		}
		if f.component != nil {
			atm.Destroy(t, f.component)
		}
	}
	err := os.RemoveAll(ts.stateDir)
	require.NoError(ts.t, err)
}

func (ts *Suite) Setup(t *testing.T, f func(t *testing.T, atm *Atmos)) {
	atm := &Atmos{}
	f(t, atm)
}

func (ts *Suite) TearDown(t *testing.T, f func(t *testing.T, atm *Atmos)) {
	ts.teardown = append(ts.teardown, &teadDown{component: nil, callback: &f})
}

func (ts *Suite) Test(t *testing.T, name string, f func(t *testing.T, atm *Atmos)) {
	if !*skipTests {
		atm := ts.getTestAtmos(ts.t)
		t.Run(name, func(t *testing.T) {
			f(t, atm)
		})
	}
}

func (ts *Suite) getAtmosOptions(vars map[string]interface{}) *atmos.Options {
	result, err := ts.options.Clone()
	require.NoError(ts.t, err)

	envvars := map[string]string{
		"TEST_SUITE_NAME": ts.name,
	}

	err = mergo.Merge(&result.EnvVars, envvars)
	require.NoError(ts.t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(ts.t, err)

	return result
}

// func (ts *Suite) GetOptions(t *testing.T, component *AtmosComponent) *atmos.Options {
// 	suiteOptions := ts.getAtmosOptions(t, &atmos.Options{}, map[string]interface{}{})
// 	return component.getAtmosOptions(t, suiteOptions, map[string]interface{}{})
// }
