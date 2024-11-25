package aws_component_helper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gruntwork-io/terratest/modules/random"
)

func setupTestSuite(ts *TestSuite) error {
	err := createStateDir(ts.TempDir)
	if err != nil {
		return err
	}

	err = createTerraformComponentsDir(ts.TempDir)
	if err != nil {
		return err
	}

	err = copyDirectoryRecursively(ts.FixturesPath, ts.TempDir)
	if err != nil {
		return err
	}

	return nil
}

func createStateDir(tempDir string) error {
	stateDir := filepath.Join(tempDir, "state")
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		err := os.MkdirAll(stateDir, 0777)
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

func readOrCreateTestSuiteFile(testSuite *TestSuite, testName string) (*TestSuite, error) {
	if data, err := os.ReadFile(testSuiteFile); err == nil {
		if err := json.Unmarshal(data, &testSuite); err != nil {
			return &TestSuite{}, fmt.Errorf("failed to parse test_suites.json: %s", err.Error())
		}

		fmt.Printf("running tests in %s\n", testSuite.TempDir)
		return testSuite, nil
	} else {
		randID := random.UniqueId()
		testSuite.RandomIdentifier = strings.ToLower(randID)

		testSuite.TempDir, err = os.MkdirTemp("", testName)
		if err != nil {
			return &TestSuite{}, err
		}
		fmt.Printf("running tests in %s\n", testSuite.TempDir)

		// Write new values to file
		data, err := json.MarshalIndent(testSuite, "", "  ")

		if err != nil {
			return &TestSuite{}, err
		}

		if err := os.WriteFile(testSuiteFile, data, 0644); err != nil {
			return &TestSuite{}, err
		}
	}

	os.Setenv("ATMOS_BASE_PATH", testSuite.TempDir)
	os.Setenv("ATMOS_CLI_CONFIG_PATH", testSuite.TempDir)
	os.Setenv("TEST_ACCOUNT_ID", testSuite.AwsAccountId)

	return testSuite, nil
}
