package aws_component_helper

import (
	"dario.cat/mergo"
	"flag"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
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

	skipTests = flag.Bool("cth.skip-tests", false, "skip tests")
)

type XTestSuites struct {
	RandomIdentifier string
	AwsAccountId     string
	AwsRegion        string
	SourceDir        string
	TempDir          string
	FixturesPath     string
	suites           map[string]*XTestSuite
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
		suites:           map[string]*XTestSuite{},
	}

	describeStacksOptions := suites.getAtmosOptions(t, &atmos.Options{}, map[string]interface{}{})
	describeStacksOptions.AtmosBasePath = filepath.Join(suites.SourceDir, suites.FixturesPath)
	describeStacksOptions.EnvVars["ATMOS_BASE_PATH"] = describeStacksOptions.AtmosBasePath
	describeStacksOptions.EnvVars["ATMOS_CLI_CONFIG_PATH"] = describeStacksOptions.AtmosBasePath

	describeStacksConfigs, err := atmos.DescribeStacksE(t, describeStacksOptions)
	require.NoError(t, err)

	for stackName, stack := range describeStacksConfigs.Stacks {
		for componentName, component := range stack.Components.Terraform {
			if component.Settings.Test == nil {
				// Skip components that are not part of tests
				continue
			}

			suiteName := component.Settings.Test.Suite
			require.NotEmptyf(t, suiteName, "settings.test.suite is required for component %s in stack %s", componentName, stackName)

			validStages := []string{"testSuiteSetUp", "testSetUp", "subjectUnderTest", "assert"}
			stage := component.Settings.Test.Stage

			require.NotEmptyf(t, stage, "settings.test.stage is required for component %s in stack %s", componentName, stackName)
			require.Containsf(t, validStages, stage, "settings.test.stage should be one of %v for component %s in stack %s", validStages, componentName, stackName)

			testName := component.Settings.Test.Test
			require.False(t, stage != "testSuiteSetUp" && testName == "", "settings.test.test is required for component %s in stack %s", componentName, stackName)

			suite := suites.GetOrCreateSuite(suiteName)

			switch stage {
			case "testSuiteSetUp":
				suite.AddSetup(componentName, stackName)
			case "testSetUp":
				suite.GetOrCreateTest(testName).AddSetup(componentName, stackName)
			case "subjectUnderTest":
				suite.GetOrCreateTest(testName).SetSubject(componentName, stackName)
			case "assert":
				suite.GetOrCreateTest(testName).AddSAssert(componentName, stackName)
			}
		}
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

func (ts *XTestSuites) Run(t *testing.T, options *atmos.Options) {
	suitesOptions := ts.getAtmosOptions(t, options, map[string]interface{}{})
	if !*skipTmpDir {
		fmt.Printf("Create TMP dir: %s \n", ts.TempDir)

		err := os.Mkdir(ts.TempDir, 0777)
		assert.NoError(t, err)
		defer os.RemoveAll(ts.TempDir)

		err = copyDirectoryRecursively(ts.SourceDir, ts.TempDir)
		assert.NoError(t, err)
	} else {
		fmt.Printf("Use source dir: %s \n", ts.SourceDir)
	}

	if !*skipVendorDependencies {
		atmosVendorPull(t, suitesOptions)
	} else {
		fmt.Println("Skip Vendor Pull")
	}

	err := createStateDir(ts.TempDir)
	assert.NoError(t, err)

	if *runParallel {
		fmt.Println("Run suites in parallel mode")
		t.Parallel()
	}
	for name, suite := range ts.suites {
		t.Run(name, func(t *testing.T) {
			suite.Run(t, suitesOptions)
		})
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

func (ts *XTestSuites) GetOrCreateSuite(name string) *XTestSuite {
	if _, ok := ts.suites[name]; !ok {
		ts.suites[name] = NewXTestSuite()
	}
	return ts.suites[name]

}
