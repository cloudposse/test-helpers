package aws_component_helper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

//func setupTestSuite(ts *TestSuite) error {
//	err := createStateDir(ts.TempDir)
//	if err != nil {
//		return err
//	}
//
//	err = createTerraformComponentsDir(ts.TempDir)
//	if err != nil {
//		return err
//	}
//
//	err = copyDirectoryRecursively(ts.FixturesPath, ts.TempDir)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func createDir(tempDir string, name string) error {
	dir := filepath.Join(tempDir, name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		return err
	}

	return nil
}

func createTerraformComponentsDir(tempDir string) error {
	stateDir := filepath.Join(tempDir, "components", "terraform")
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		err := os.MkdirAll(stateDir, 0777)
		return err
	}

	return nil
}

func getTestName() (string, error) {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return "", fmt.Errorf("unable to get the test file name")
	}

	testName := filepath.Base(filename[:len(filename)-3]) + "_"
	return testName, nil
}

func getAwsAccountId() (string, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}
	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return *identity.Account, nil
}

//func readOrCreateTestSuiteFile(testSuite *TestSuite, testName string) (*TestSuite, error) {
//	// Initialize TestSuites structure
//	var testSuites TestSuites
//
//	if data, err := os.ReadFile(testSuiteFile); err == nil {
//		// File exists, try to unmarshal existing test suites
//		if err := json.Unmarshal(data, &testSuites); err != nil {
//			return &TestSuite{}, fmt.Errorf("failed to parse test_suites.json: %s", err.Error())
//		}
//
//		if len(testSuites.Suites) > 1 && testSuite.Index < 0 {
//			return &TestSuite{}, fmt.Errorf("test suite index is required when multiple test suites are present")
//		}
//		if testSuite.Index == -1 && len(testSuites.Suites) == 1 {
//			testSuite.Index = 0
//		}
//		if !testSuite.ForceNewSuite && len(testSuites.Suites) > 0 {
//			return testSuites.Suites[testSuite.Index], nil
//		}
//	}
//
//	// If we get here, either the file doesn't exist or we didn't find a matching suite
//	fmt.Println("no matching test suite found for index", testSuite.Index, "creating new test suite")
//	randID := random.UniqueId()
//	testSuite.RandomIdentifier = strings.ToLower(randID)
//	testSuite.Index = len(testSuites.Suites) // Set index to current length
//
//	var err error
//	testSuite.TempDir, err = os.MkdirTemp("", testName)
//	if err != nil {
//		return &TestSuite{}, err
//	}
//	fmt.Printf("running tests in %s\n", testSuite.TempDir)
//
//	// Add new test suite to the collection
//	testSuites.Suites = append(testSuites.Suites, testSuite)
//
//	// Write updated test suites to file
//	data, err := json.MarshalIndent(testSuites, "", "  ")
//	if err != nil {
//		return &TestSuite{}, err
//	}
//
//	if err := os.WriteFile(testSuiteFile, data, 0644); err != nil {
//		return &TestSuite{}, err
//	}
//
//	// os.Setenv("ATMOS_BASE_PATH", testSuite.TempDir)
//	os.Setenv("ATMOS_CLI_CONFIG_PATH", testSuite.TempDir)
//	os.Setenv("TEST_ACCOUNT_ID", testSuite.AwsAccountId)
//
//	return testSuite, nil
//}
