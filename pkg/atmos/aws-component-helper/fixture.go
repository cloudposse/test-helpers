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

// Define global callback function to retrieve AWS account ID
var (
	getAwsAaccountIdCallback = getAwsAccountId
)

// Define global flags for configuring test behavior
var (
	skipTmpDir             = flag.Bool("skip-tmp-dir", false, "Run in the current directory")
	skipVendorDependencies = flag.Bool("skip-vendor", false, "skip vendor dependencies")
	skipTeardownFixtures   = flag.Bool("skip-fixtures-teardown", false, "skip fixtures teardown")
	skipTeardown           = flag.Bool("skip-teardown", false, "skip teardown")
	useCache               = flag.Bool("cache", false, "use cache for terraform plugins")
	matchSuiteAndTest      = flag.String("match", "", "regular expression to select suite and tests to run")
)

// Fixture struct holds test-specific configurations and state
type Fixture struct {
	t                *testing.T // Testing object
	RandomIdentifier string     // Unique identifier for the fixture
	AwsAccountId     string     // AWS Account ID
	AwsRegion        string     // AWS Region
	SourceDir        string     // Source directory for the fixture
	TempDir          string     // Temporary directory for the fixture
	FixturesPath     string     // Path to the fixture directory
	State            *State     // State management object
	suites           []*Suite   // List of test suites
	suitesNames      []string   // Names of test suites
}

// NewFixture initializes a new Fixture instance
func NewFixture(t *testing.T, sourceDir string, awsRegion string, fixturesDir string) *Fixture {
	awsAccountId, err := getAwsAaccountIdCallback()
	require.NoError(t, err) // Ensure AWS account ID retrieval succeeded

	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	// Create a temporary directory path for the fixture
	tmpdir := filepath.Join(os.TempDir(), "fixtures-"+randomId)

	realSourcePath, err := filepath.Abs(sourceDir)
	require.NoError(t, err) // Ensure source directory resolution succeeded

	// Initialize state with namespace and state directory
	state_namespace := strings.ReplaceAll(filepath.Join(t.Name(), fixturesDir), "/", "-")
	state := NewState(state_namespace, filepath.Join(realSourcePath, "state"))

	// Return initialized Fixture instance
	return &Fixture{
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
}

// WorkDir returns the working directory based on the skipTmpDir flag
func (ts *Fixture) WorkDir() string {
	if !*skipTmpDir {
		return ts.TempDir
	} else {
		return ts.SourceDir
	}
}

// FixtureDir returns the path to the fixture directory
func (ts *Fixture) FixtureDir() string {
	return filepath.Join(ts.WorkDir(), ts.FixturesPath)
}

// SetUp prepares the fixture for use
func (ts *Fixture) SetUp(options *atmos.Options) {
	suitesOptions := ts.getAtmosOptions(options, map[string]interface{}{})

	// Create temporary directory if skipTmpDir flag is not set
	if !*skipTmpDir {
		fmt.Printf("Create fixtures tmp dir: %s \n", ts.TempDir)

		err := os.Mkdir(ts.TempDir, 0777)
		require.NoError(ts.t, err) // Ensure directory creation succeeded

		err = copyDirectoryRecursively(ts.SourceDir, ts.TempDir)
		require.NoError(ts.t, err) // Ensure directory copy succeeded
	} else {
		fmt.Printf("Use fixtures source dir: %s \n", ts.SourceDir)
	}

	// Handle vendor dependencies based on skipVendorDependencies flag
	if !*skipVendorDependencies {
		atmosVendorPull(ts.t, suitesOptions)
	} else {
		fmt.Println("Skip Vendor Pull")
	}

	// Set up state
	err := ts.State.SetUp()
	require.NoError(ts.t, err)

	// Create cache directory if useCache flag is set
	if *useCache {
		err = createDir(ts.WorkDir(), ".cache")
		require.NoError(ts.t, err)
	}
}

// TearDown cleans up resources created by the fixture
func (ts *Fixture) TearDown() {
	if *skipTeardown {
		fmt.Println("Skip teardown")
		return
	}

	// Run teardown for all suites
	for i := len(ts.suites) - 1; i >= 0; i-- {
		ts.suites[i].runTeardown()
	}

	// Remove temporary directory if skipTmpDir and skipTeardownFixtures flags are not set
	if !*skipTmpDir && !*skipTeardownFixtures {
		err := os.RemoveAll(ts.TempDir)
		require.NoError(ts.t, err)
	}

	// Tear down state
	err := ts.State.Teardown()
	require.NoError(ts.t, err)
}

// getAtmosOptions generates options for Atmos CLI with environment variables and state
func (ts *Fixture) getAtmosOptions(options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone() // Clone existing options
	}

	// Configure Atmos options
	result.AtmosBasePath = ts.FixtureDir()
	result.NoColor = true
	result.Lock = false
	result.Upgrade = true

	// Define environment variables
	envvars := map[string]string{
		"TEST_ACCOUNT_ID":       ts.AwsAccountId,
		"ATMOS_BASE_PATH":       result.AtmosBasePath,
		"ATMOS_CLI_CONFIG_PATH": result.AtmosBasePath,
	}

	// Add cache directory if useCache flag is set
	if *useCache {
		envvars["TF_PLUGIN_CACHE_DIR"] = filepath.Join(ts.WorkDir(), ".cache")
	}

	err := mergo.Merge(&result.EnvVars, envvars) // Merge environment variables
	require.NoError(ts.t, err)

	// Define and merge fixture-specific variables
	fixtureVars := map[string]interface{}{
		"region": ts.AwsRegion,
	}

	err = mergo.Merge(&result.Vars, fixtureVars)
	require.NoError(ts.t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(ts.t, err)

	return result
}

// Suite adds a test suite to the Fixture
func (ts *Fixture) Suite(name string, f func(t *testing.T, suite *Suite)) {
	require.NotContains(ts.t, ts.suitesNames, name, "Suite %s already exists", name)

	suite := NewSuite(ts.t, name, ts) // Create new suite

	ts.suites = append(ts.suites, suite)
	ts.suitesNames = append(ts.suitesNames, name)

	// Run the suite if it matches the filter
	suiteRunName := fmt.Sprintf("%s/%s", ts.t.Name(), suite.name)
	if ok, err := matchFilter(suiteRunName); ok {
		ts.t.Run(name, func(t *testing.T) {
			f(t, suite)
		})
	} else {
		require.NoError(ts.t, err)
	}
}
