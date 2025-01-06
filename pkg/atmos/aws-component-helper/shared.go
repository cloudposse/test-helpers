package aws_component_helper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// copyDirectoryRecursively copies all contents of the source directory (srcDir) to the destination directory (destDir).
func copyDirectoryRecursively(srcDir string, destDir string) error {
	fmt.Println("copying contents of", srcDir, "to", destDir)

	// Walk through all files and directories in the source directory
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path from srcDir
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Create the corresponding destination path
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create the directory at the destination
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// Read the file contents
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Write the file contents to the destination
			return os.WriteFile(destPath, content, info.Mode())
		}
	})
}

// parseListOutputTerraform converts a list of interface{} (as returned by Terraform outputs) to a list of strings.
func parseListOutputTerraform(outputList []interface{}) ([]string, error) {
	list := []string{}

	// Convert each item to a string
	for _, item := range outputList {
		list = append(list, fmt.Sprintf("%v", item))
	}

	return list, nil
}

// matchFilter checks if a test name matches a filter (defined in *matchSuiteAndTest).
// The filter is split into parts and matched against the name using regex.
func matchFilter(name string) (bool, error) {
	nameParts := strings.Split(name, "/")

	// Validate that the name has at least one '/'
	if len(nameParts) == 1 {
		return false, fmt.Errorf("Invalid test name: %s. Should contains at least 1 '/'", name)
	}
	nameParts = nameParts[1:] // Exclude the first part of the name

	// Split the filter into parts
	matchParts := strings.Split(*matchSuiteAndTest, "/")

	// Match only up to the smallest number of parts
	partsCount := min(len(nameParts), len(matchParts))

	for i := 0; i < partsCount; i++ {
		result, err := regexp.MatchString(matchParts[i], nameParts[i])
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}

	return true, nil
}

// createDir creates a directory (and any necessary parent directories) at the specified location.
func createDir(tempDir string, name string) error {
	dir := filepath.Join(tempDir, name)

	// Check if the directory exists; create it if it doesn't
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		return err
	}

	return nil
}

// createTerraformComponentsDir creates the directory structure for Terraform components under the given tempDir.
func createTerraformComponentsDir(tempDir string) error {
	stateDir := filepath.Join(tempDir, "components", "terraform")

	// Check if the directory exists; create it if it doesn't
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		err := os.MkdirAll(stateDir, 0777)
		return err
	}

	return nil
}

// getTestName retrieves the name of the current test file (without the extension) and appends an underscore.
func getTestName() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get the test file name")
	}

	// Remove the file extension and append an underscore
	testName := filepath.Base(filename[:len(filename)-3]) + "_"
	return testName, nil
}

// getAwsAccountId retrieves the AWS account ID of the caller using the AWS STS service.
func getAwsAccountId() (string, error) {
	ctx := context.Background()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	// Configure retries for reliability
	cfg.RetryMode = aws.RetryModeAdaptive
	cfg.RetryMaxAttempts = 3

	// Create an STS client
	stsClient := sts.NewFromConfig(cfg)

	// Get the caller's identity
	identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	// Return the AWS account ID
	return *identity.Account, nil
}
