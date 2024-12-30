package aws_component_helper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func copyDirectoryRecursively(srcDir string, destDir string) error {
	fmt.Println("copying contents of", srcDir, "to", destDir)
	// Walk through all files and directories in srcDir and copy them to destDir
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path from srcDir
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Calculate destination path in destDir
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory in destination
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// Copy file content
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			return os.WriteFile(destPath, content, info.Mode())
		}
	})
}

// Parse a list output in the format it is returned by Terraform 0.12 and newer versions
func parseListOutputTerraform(outputList []interface{}, key string) ([]string, error) {
	list := []string{}

	for _, item := range outputList {
		list = append(list, fmt.Sprintf("%v", item))
	}

	return list, nil
}

func matchFilter(name string) (bool, error) {
	nameParts := strings.Split(name, "/")

	if len(nameParts) == 1 {
		return false, fmt.Errorf("Invalid test name: %s. Should contains at least 1 '/'", name)
	}
	nameParts = nameParts[1:]

	matchParts := strings.Split(*matchSuiteAndTest, "/")

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
