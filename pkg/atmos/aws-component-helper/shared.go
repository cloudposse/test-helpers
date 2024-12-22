package aws_component_helper

import (
	"fmt"
	"os"
	"path/filepath"
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
