package component

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

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

func setupDependencyStacks(m *testing.M, suite *TestSuite, stacks []Stack) error {
	fmt.Printf("Setting up dependency stacks in %s\n", suite.TempDir)
	for _, stack := range stacks {
		fmt.Printf("    Component: %s, Stack: %s\n", stack.Component, stack.StackName)
	}

	return nil
}

func tearDownDependencyStacks(m *testing.M, suite *TestSuite, stacks []Stack) error {
	// iterate through all the stacks and print the component name and stack name
	for _, stack := range stacks {
		fmt.Printf("Tearing down dependency Component: %s, Stack: %s\n", stack.Component, stack.StackName)
	}

	return nil
}

func setupTest(m *testing.M) error {
	fmt.Println("Setting up test")
	return nil
}

func tearDownTest(m *testing.M) error {
	fmt.Println("Tearing down test")
	return nil
}

func AtmosAwsComponentMainHelper(m *testing.M, opts AwsComponentTestOptions) {
	cliArgs := parseCLIArgs()

	testSuite, err := setupTestSuite()
	if err != nil {
		panic("Failed to setup test suite: " + err.Error())
	}

	if !cliArgs.SkipDependencyTeardown {
		defer tearDownDependencyStacks(m, &testSuite, opts.StackDependencies)
	}

	if !cliArgs.SkipDependencySetup {
		setupDependencyStacks(m, &testSuite, opts.StackDependencies)
	}

	if !cliArgs.SkipTestTeardown {
		defer tearDownTest(m)
	}

	if !cliArgs.SkipTestSetup {
		setupTest(m)
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
