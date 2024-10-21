package component

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/cloudposse/test-helpers/pkg/awsnuke"
	tt "github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

type AtmosComponentCLIOptions struct {
	SkipDependencySetup    bool
	SkipDependencyTeardown bool
	SkipTestSetup          bool
	SkipTestTeardown       bool
}

type TestSuite struct {
	TempDir    string
	RandomSeed string
}

func parseCLIArgs() *AtmosComponentCLIOptions {
	skipDependencySetup := flag.Bool("skip-deps-setup", false, "skip dependency setup")
	skipDependencyTeardown := flag.Bool("skip-deps-teardown", false, "skip dependency teardown")
	skipTestSetup := flag.Bool("skip-test-setup", false, "skip test setup")
	skipTestTeardown := flag.Bool("skip-test-teardown", false, "skip test teardown")
	flag.Parse()

	return &AtmosComponentCLIOptions{
		SkipDependencySetup:    *skipDependencySetup,
		SkipDependencyTeardown: *skipDependencyTeardown,
		SkipTestSetup:          *skipTestSetup,
		SkipTestTeardown:       *skipTestTeardown,
	}
}

func setupTestSuite() (TestSuite, error) {
	var testSuite TestSuite
	const testSuiteFile = "test_suite.json"

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get the test file name")
	}
	testName := filepath.Base(filename[:len(filename)-3]) + "_"

	if data, err := os.ReadFile(testSuiteFile); err == nil {
		if err := json.Unmarshal(data, &testSuite); err != nil {
			panic("Failed to parse test_suites.json: " + err.Error())
		}

		return testSuite, nil
	} else {
		randID := strconv.Itoa(rand.Intn(100000))
		testSuite.RandomSeed = randID

		testSuite.TempDir, err = os.MkdirTemp("", testName)
		if err != nil {
			return TestSuite{}, err
		}

		// Write new values to file
		data, err := json.MarshalIndent(testSuite, "", "  ")

		if err != nil {
			return TestSuite{}, err
		}

		if err := os.WriteFile(testSuiteFile, data, 0644); err != nil {
			return TestSuite{}, err
		}
	}

	fmt.Println("Setting up test suites")
	return testSuite, nil
}

func tearDownTestSuite() error {
	fmt.Println("Tearing down test suites")
	return nil
}

func setupDependencyStacks(t tt.TestingT, suite *TestSuite, opts AwsComponentTestOptions) {
	fmt.Printf("Setting up dependency stacks in %s\n", suite.TempDir)
	for _, dependency := range opts.StackDependencies {
		fmt.Printf("Applying dependency: %s in stack %s\n", dependency.Component, dependency.StackName)
		_, err := atmos.ApplyE(t, atmos.WithDefaultRetryableErrors(t, &atmos.Options{
			AtmosBasePath: suite.TempDir,
			Component:     dependency.Component,
			Stack:         dependency.StackName,
			NoColor:       true,
			BackendConfig: map[string]interface{}{
				"workspace_key_prefix": strings.Join([]string{suite.RandomSeed, opts.StackName}, "-"),
			},
			Vars: map[string]interface{}{
				"attributes": []string{suite.RandomSeed},
				"default_tags": map[string]string{
					"CreatedByTerratestRun": suite.RandomSeed,
				},
				"region": dependency.Region,
			},
		}))
		require.NoError(t, err)
	}
}

func tearDownDependencyStacks(t tt.TestingT, suite *TestSuite, opts AwsComponentTestOptions) {
	for _, dependency := range opts.StackDependencies {
		t.Logf("Destroying dependency: %s in stack %s\n", dependency.Component, dependency.StackName)
		_, err := atmos.DestroyE(t, atmos.WithDefaultRetryableErrors(t, &atmos.Options{
			AtmosBasePath: suite.TempDir,
			Component:     dependency.Component,
			Stack:         dependency.StackName,
			NoColor:       true,
			BackendConfig: map[string]interface{}{
				"workspace_key_prefix": strings.Join([]string{suite.RandomSeed, opts.StackName}, "-"),
			},
			Vars: map[string]interface{}{
				"attributes": []string{suite.RandomSeed},
				"default_tags": map[string]string{
					"CreatedByTerratestRun": suite.RandomSeed,
				},
				"region": dependency.Region,
			},
		}))
		require.NoError(t, err)
	}
}

func setupTest(t tt.TestingT, suite *TestSuite, opts AwsComponentTestOptions) string {

	t.Log("Performing test setup...")

	atmosOptions := atmos.WithDefaultRetryableErrors(t, &atmos.Options{
		AtmosBasePath: suite.TempDir,
		Component:     opts.ComponentName,
		Stack:         opts.StackName,
		NoColor:       true,
		BackendConfig: map[string]interface{}{
			"workspace_key_prefix": strings.Join([]string{suite.RandomSeed, opts.StackName}, "-"),
		},
		Vars: map[string]interface{}{
			"attributes": []string{suite.RandomSeed},
			"default_tags": map[string]string{
				"CreatedByTerratestRun": suite.RandomSeed,
			},
			"region": opts.AwsRegion,
		},
	})
	options := atmos.WithDefaultRetryableErrors(t, atmosOptions)
	out := atmos.Apply(t, options)

	return out
}

func tearDownTest(m *testing.M) error {
	fmt.Println("Tearing down test")
	return nil
}

func AtmosAwsComponentMainHelper(m *testing.M, opts AwsComponentTestOptions) {
	// Create a mocked testing.T instance so we can call methods that expect a testing.T instance from TestMain
	t := &tt.CustomT{}

	cliArgs := parseCLIArgs()

	testSuite, err := setupTestSuite()
	if err != nil {
		panic("Failed to setup test suite: " + err.Error())
	}

	if !cliArgs.SkipDependencyTeardown {
		defer tearDownDependencyStacks(t, &testSuite, opts)
	}

	if !cliArgs.SkipDependencySetup {
		setupDependencyStacks(t, &testSuite, opts)
	}

	if !cliArgs.SkipTestTeardown {
		defer tearDownTest(m)
	}

	if !cliArgs.SkipTestSetup {
		setupTest(t, &testSuite, opts)
	}

	if !opts.SkipAwsNuke {
		awsnuke.NukeTestAccountByTag(t, "CreatedByTerratestRun", testSuite.RandomSeed, []string{opts.AwsRegion}, false)
	}

	m.Run()
}

// func AtmosAwsComponentTestHelper(t *testing.T, opts AwsComponentTestOptions, callback func(t *testing.T, opts *atmos.Options, output string)) {
// 	t.Helper()

// 	// Setup temp dir
// 	// Copy src to temp dir/component/component_name
// 	// Run atmos vendor pull
// 	// Terraform apply dependencies
// 	// Terraform apply component
// 	// Run callback
// 	// Terraform destroy component
// 	// Terraform destroy dependencies
// 	// Delete temp dir

// 	output := ""
// 	callback(t, opts, output)
// }
