package aws_component_helper

import (
	"flag"
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

var (
	skipTests      = flag.Bool("skip-tests", false, "skip tests")
	skipSetup      = flag.Bool("skip-setup", false, "skip setup")
	preserveStates = flag.Bool("preserve-states", true, "preserve states")

	// runParallel            = flag.Bool("parallel", false, "Run parallel")

	// forceNewSuite           = flag.Bool("cth.force-new-suite", false, "force new suite")
	// suiteIndex              = flag.Int("cth.suite-index", -1, "suite index")
	// skipAwsNuke             = flag.Bool("cth.skip-aws-nuke", false, "skip aws nuke")

	// skipDependencies  = flag.Bool("cth.skip-deps", false, "skip deploy dependencies")
	// skipDeployDependencies  = flag.Bool("cth.skip-deps-deploy", false, "skip deploy dependencies")
	// skipDestroyDependencies = flag.Bool("cth.skip-deps-destroy", false, "skip destroy dependencies")
	// skipTeardownTestSuite = flag.Bool("skip-teardown", false, "skip test suite teardown")

	// skipDeployComponentUnderTest  = flag.Bool("cth.skip-deploy-cut", false, "skip deploy component under test")
	// skipDestroyComponentUnderTest = flag.Bool("cth.skip-destroy-cut", false, "skip destroy component under test")
)

type teadDown struct {
	component *AtmosComponent
	callback  *func(t *testing.T, atm *Atmos)
}

type Suite struct {
	t                *testing.T
	randomIdentifier string
	name             string
	globalStateDir   string
	dependencies     []*AtmosComponent
	teardown         []*teadDown
	options          *atmos.Options
}

func NewSuite(t *testing.T, name string, fixture *Fixture) *Suite {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	require.NotContains(t, name, "/")
	suite := &Suite{
		t:                t,
		name:             name,
		randomIdentifier: randomId,
		dependencies:     []*AtmosComponent{},
		globalStateDir:   fixture.GlobalStateDir(),
		teardown:         []*teadDown{},
		options:          fixture.getAtmosOptions(&atmos.Options{}, map[string]interface{}{}),
	}

	if fixture.StateDir() != "" {
		copyDirectoryRecursively(fixture.StateDir(), suite.StateDir())
	} else {
		err := createDir(fixture.GlobalStateDir(), suite.stateNamespace())
		require.NoError(t, err)
	}

	return suite
}

func (ts *Suite) StateDir() string {
	return filepath.Join(ts.globalStateDir, ts.stateNamespace())
}

func (ts *Suite) stateNamespace() string {
	return fmt.Sprintf("suite-%s", ts.name)

}

func (ts *Suite) AddDependency(componentName string, stackName string) {
	if *skipSetup {
		fmt.Printf("Skip suite %s setup dependency component: %s stack: %s\n", ts.name, componentName, stackName)
		return
	}
	component := NewAtmosComponent(componentName, stackName, nil)
	ts.dependencies = append(ts.dependencies, component)
	ts.teardown = append(ts.teardown, &teadDown{component: component, callback: nil})
	ts.getAtmos().Deploy(component)
}

func (ts *Suite) getAtmos() *Atmos {
	return NewAtmos(ts.t, ts.getAtmosOptions(map[string]interface{}{}))
}

func (ts *Suite) getTestAtmos() *Atmos {
	return ts.getAtmos()
}

func (ts *Suite) runTeardown() {
	if *skipTeardown {
		fmt.Printf("Skip teardown suite %s\n", ts.name)
		fmt.Printf("Suite %s preserve states %s\n", ts.name, ts.StateDir())
		return
	}
	atm := ts.getAtmos()
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
	if *preserveStates {
		fmt.Printf("Suite %s preserve states %s\n", ts.name, ts.StateDir())
	} else {
		fmt.Printf("Suite %s drops states %s\n", ts.name, ts.StateDir())
		err := os.RemoveAll(ts.StateDir())
		require.NoError(ts.t, err)
	}
}

func (ts *Suite) Setup(t *testing.T, f func(t *testing.T, atm *Atmos)) {
	if *skipSetup {
		fmt.Printf("Skip suite %s setup callback\n", ts.name)
		return
	}
	atm := ts.getAtmos()
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
	atm := ts.getTestAtmos()
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

	envvars := map[string]string{
		"TEST_SUITE_NAME": ts.stateNamespace(),
	}

	err = mergo.Merge(&result.EnvVars, envvars)
	require.NoError(ts.t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(ts.t, err)

	return result
}
