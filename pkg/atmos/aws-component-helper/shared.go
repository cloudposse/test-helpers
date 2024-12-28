package aws_component_helper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
		fmt.Printf("Matching %s with %s\n", matchParts[i], nameParts[i])
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
