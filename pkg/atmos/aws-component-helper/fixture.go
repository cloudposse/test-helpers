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
	getAwsAaccountIdCallback = getAwsAccountId
)

var (
	skipTmpDir             = flag.Bool("skip-tmp-dir", false, "Run in the current directory")
	skipVendorDependencies = flag.Bool("skip-vendor", false, "skip vendor dependencies")
	skipTeardownFixtures   = flag.Bool("skip-fixtures-teardown", false, "skip fixtures teardown")
	skipTeardown           = flag.Bool("skip-teardown", false, "skip teardown")
	useCache               = flag.Bool("cache", false, "use cache for terraform plugins")
	matchSuiteAndTest      = flag.String("match", "", "regular expression to select suite and tests to run")
)

type Fixture struct {
	t                *testing.T
	RandomIdentifier string
	AwsAccountId     string
	AwsRegion        string
	SourceDir        string
	TempDir          string
	FixturesPath     string
	State            *State
	suites           []*Suite
	suitesNames      []string
}

func NewFixture(t *testing.T, sourceDir string, awsRegion string, fixturesDir string) *Fixture {
	awsAccountId, err := getAwsAaccountIdCallback()
	require.NoError(t, err)

	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	tmpdir := filepath.Join(os.TempDir(), "fixtures-"+randomId)

	realSourcePath, err := filepath.Abs(sourceDir)
	require.NoError(t, err)

	state_namespace := strings.ReplaceAll(filepath.Join(t.Name(), fixturesDir), "/", "-")
	state := NewState(state_namespace, filepath.Join(realSourcePath, "state"))

	suites := &Fixture{
		t:                t,
		RandomIdentifier: randomId,
		SourceDir:        realSourcePath,
		TempDir:          tmpdir,
		FixturesPath:     fixturesDir,
		AwsAccountId:     awsAccountId,
		AwsRegion:        awsRegion,
		suites:           []*Suite{},
		suitesNames:      []string{},
		State:            state,
	}

	return suites
}

func (ts *Fixture) WorkDir() string {
	if !*skipTmpDir {
		return ts.TempDir
	} else {
		return ts.SourceDir
	}
}

func (ts *Fixture) FixtureDir() string {
	return filepath.Join(ts.WorkDir(), ts.FixturesPath)
}

func (ts *Fixture) SetUp(options *atmos.Options) {
	suitesOptions := ts.getAtmosOptions(options, map[string]interface{}{})
	if !*skipTmpDir {
		fmt.Printf("Create fixtures tmp dir: %s \n", ts.TempDir)

		err := os.Mkdir(ts.TempDir, 0777)
		require.NoError(ts.t, err)

		err = copyDirectoryRecursively(ts.SourceDir, ts.TempDir)
		require.NoError(ts.t, err)
	} else {
		fmt.Printf("Use fixtures source dir: %s \n", ts.SourceDir)
	}

	if !*skipVendorDependencies {
		atmosVendorPull(ts.t, suitesOptions)
	} else {
		fmt.Println("Skip Vendor Pull")
	}

	err := ts.State.SetUp()
	require.NoError(ts.t, err)

	if *useCache {
		err = createDir(ts.WorkDir(), ".cache")
		require.NoError(ts.t, err)
	}
}

func (ts *Fixture) TearDown() {
	if *skipTeardown {
		fmt.Println("Skip teardown")
		return
	}
	for i := len(ts.suites) - 1; i >= 0; i-- {
		ts.suites[i].runTeardown()
	}
	if !*skipTmpDir && !*skipTeardownFixtures {
		err := os.RemoveAll(ts.TempDir)
		require.NoError(ts.t, err)
	}
	err := ts.State.Teardown()
	require.NoError(ts.t, err)
}

func (ts *Fixture) getAtmosOptions(options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone()
	}

	result.AtmosBasePath = ts.FixtureDir()
	result.NoColor = true
	result.Lock = false
	result.Upgrade = true

	envvars := map[string]string{
		"TEST_ACCOUNT_ID":       ts.AwsAccountId,
		"ATMOS_BASE_PATH":       result.AtmosBasePath,
		"ATMOS_CLI_CONFIG_PATH": result.AtmosBasePath,
	}

	if *useCache {
		envvars["TF_PLUGIN_CACHE_DIR"] = filepath.Join(ts.WorkDir(), ".cache")
	}

	err := mergo.Merge(&result.EnvVars, envvars)
	require.NoError(ts.t, err)

	fixtureVars := map[string]interface{}{
		"region": ts.AwsRegion,
	}

	err = mergo.Merge(&result.Vars, fixtureVars)
	require.NoError(ts.t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(ts.t, err)

	return result
}

func (ts *Fixture) Suite(name string, f func(t *testing.T, suite *Suite)) {
	require.NotContains(ts.t, ts.suitesNames, name, "Suite %s already exists", name)

	suite := NewSuite(ts.t, name, ts)

	ts.suites = append(ts.suites, suite)
	ts.suitesNames = append(ts.suitesNames, name)

	suiteRunName := fmt.Sprintf("%s/%s", ts.t.Name(), suite.name)
	if ok, err := matchFilter(suiteRunName); ok {
		ts.t.Run(name, func(t *testing.T) {
			f(t, suite)
		})
	} else {
		require.NoError(ts.t, err)
	}
}
