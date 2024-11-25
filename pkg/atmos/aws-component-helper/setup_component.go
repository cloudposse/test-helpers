package aws_component_helper

import (
	"fmt"
	"os"
	"path/filepath"
)

func setupComponentUnderTest(ts *TestSuite) error {
	err := createComponentDir(ts)
	if err != nil {
		return err
	}

	err = copyComponentSrcToTempDir(ts)
	if err != nil {
		return err
	}
	return nil
}

func copyComponentSrcToTempDir(ts *TestSuite) error {

	destDir := filepath.Join(ts.TempDir, "components", "terraform", ts.ComponentName)
	fmt.Println("copying contents of", ts.ComponentSrcPath, "to", destDir)
	err := copyDirectoryRecursively(ts.ComponentSrcPath, destDir)
	if err != nil {
		return err
	}

	return nil
}

func createComponentDir(ts *TestSuite) error {
	destDir := filepath.Join(ts.TempDir, "components", "terraform", ts.ComponentName)
	err := os.MkdirAll(destDir, 0777)
	if err != nil {
		return err
	}
	return nil
}
