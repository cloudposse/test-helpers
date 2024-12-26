package aws_component_helper

import (
	"dario.cat/mergo"
	"flag"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	getAwsAaccountIdCallback = getAwsAccountId
)

var (
	skipTmpDir              = flag.Bool("cth.skip-tmp-dir", false, "Run in the current directory")
	skipVendorDependencies  = flag.Bool("cth.skip-vendor", false, "skip vendor dependencies")
	runParallel             = flag.Bool("cth.parallel", false, "Run parallel")
	forceNewSuite           = flag.Bool("cth.force-new-suite", false, "force new suite")
	suiteIndex              = flag.Int("cth.suite-index", -1, "suite index")
	skipAwsNuke             = flag.Bool("cth.skip-aws-nuke", false, "skip aws nuke")
	skipDeployDependencies  = flag.Bool("cth.skip-deploy-deps", false, "skip deploy dependencies")
	skipDestroyDependencies = flag.Bool("cth.skip-destroy-deps", false, "skip destroy dependencies")
	skipTeardownTestSuite   = flag.Bool("cth.skip-teardown", false, "skip test suite teardown")
	skipTests               = flag.Bool("cth.skip-tests", false, "skip tests")

	skipDeployComponentUnderTest  = flag.Bool("cth.skip-deploy-cut", false, "skip deploy component under test")
	skipDestroyComponentUnderTest = flag.Bool("cth.skip-destroy-cut", false, "skip destroy component under test")
)

type XTestSuites struct {
	RandomIdentifier string
	AwsAccountId     string
	AwsRegion        string
	SourceDir        string
	TempDir          string
	FixturesPath     string
}

func NewTestSuites(t *testing.T, sourceDir string, awsRegion string, fixturesDir string) *XTestSuites {
	awsAccountId, err := getAwsAaccountIdCallback()
	require.NoError(t, err)

	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	tmpdir := filepath.Join(os.TempDir(), "test-suites-"+randomId)
	realSourcePath, err := filepath.Abs(sourceDir)
	require.NoError(t, err)
	suites := &XTestSuites{
		RandomIdentifier: randomId,
		SourceDir:        realSourcePath,
		TempDir:          tmpdir,
		FixturesPath:     fixturesDir,
		AwsAccountId:     awsAccountId,
		AwsRegion:        awsRegion,
	}

	return suites
}

func (ts *XTestSuites) WorkDir() string {
	if !*skipTmpDir {
		return filepath.Join(ts.TempDir, ts.FixturesPath)
	} else {
		return filepath.Join(ts.SourceDir, ts.FixturesPath)
	}
}

func (ts *XTestSuites) SetUp(t *testing.T, options *atmos.Options) {
	suitesOptions := ts.getAtmosOptions(t, options, map[string]interface{}{})
	if !*skipTmpDir {
		fmt.Printf("Create TMP dir: %s \n", ts.TempDir)

		err := os.Mkdir(ts.TempDir, 0777)
		require.NoError(t, err)

		err = copyDirectoryRecursively(ts.SourceDir, ts.TempDir)
		require.NoError(t, err)
	} else {
		fmt.Printf("Use source dir: %s \n", ts.SourceDir)
	}

	if !*skipVendorDependencies {
		atmosVendorPull(t, suitesOptions)
	} else {
		fmt.Println("Skip Vendor Pull")
	}

	if !*skipTmpDir {
		err := createStateDir(ts.TempDir)
		require.NoError(t, err)
	} else {
		err := createStateDir(ts.SourceDir)
		require.NoError(t, err)
	}
}

func (ts *XTestSuites) TearDown(t *testing.T) {
	if !*skipTmpDir {
		err := os.RemoveAll(ts.TempDir)
		require.NoError(t, err)
	}
}

func (ts *XTestSuites) getAtmosOptions(t *testing.T, options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone()
	}

	result.AtmosBasePath = ts.WorkDir()
	result.NoColor = true
	result.Lock = false
	result.Upgrade = true

	envvars := map[string]string{
		"TEST_ACCOUNT_ID":       ts.AwsAccountId,
		"ATMOS_BASE_PATH":       result.AtmosBasePath,
		"ATMOS_CLI_CONFIG_PATH": result.AtmosBasePath,
	}

	err := mergo.Merge(&result.EnvVars, envvars)
	require.NoError(t, err)

	suiteVars := map[string]interface{}{
		"region": ts.AwsRegion,
	}

	err = mergo.Merge(&result.Vars, suiteVars)
	require.NoError(t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(t, err)

	return result
}

func (ts *XTestSuites) DeployComponent(t *testing.T, component *AtmosComponent, options *atmos.Options) {
	if !*skipDeployComponentUnderTest {
		suiteOptions := ts.getAtmosOptions(t, options, map[string]interface{}{})
		componentOptions := component.getAtmosOptions(t, suiteOptions, map[string]interface{}{})
		atmosApply(t, componentOptions)
	}
}

func (ts *XTestSuites) DestroyComponent(t *testing.T, component *AtmosComponent, options *atmos.Options) {
	if !*skipDeployComponentUnderTest && !*skipDestroyComponentUnderTest {
		suiteOptions := ts.getAtmosOptions(t, options, map[string]interface{}{})
		componentOptions := component.getAtmosOptions(t, suiteOptions, map[string]interface{}{})
		atmosDestroy(t, componentOptions)
	}
}

func (ts *XTestSuites) DeployDependency(t *testing.T, component *AtmosComponent, options *atmos.Options) {
	if !*skipDeployDependencies {
		suiteOptions := ts.getAtmosOptions(t, options, map[string]interface{}{})
		componentOptions := component.getAtmosOptions(t, suiteOptions, map[string]interface{}{})
		atmosApply(t, componentOptions)
	}
}

func (ts *XTestSuites) DestroyDependency(t *testing.T, component *AtmosComponent, options *atmos.Options) {
	if !*skipDeployDependencies && !*skipDestroyDependencies {
		suiteOptions := ts.getAtmosOptions(t, options, map[string]interface{}{})

		componentOptions := component.getAtmosOptions(t, suiteOptions, map[string]interface{}{})
		atmosDestroy(t, componentOptions)
	}
}

func (ts *XTestSuites) CreateAndDeployDependency(t *testing.T, componentName string, stackName string, options *atmos.Options) *AtmosComponent {
	component := NewAtmosComponent(componentName, stackName)
	ts.DeployDependency(t, component, options)
	return component
}

func (ts *XTestSuites) CreateAndDeployComponent(t *testing.T, componentName string, stackName string, options *atmos.Options) *AtmosComponent {
	component := NewAtmosComponent(componentName, stackName)
	ts.DeployComponent(t, component, options)
	return component
}
