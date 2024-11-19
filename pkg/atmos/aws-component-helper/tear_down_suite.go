package aws_component_helper

import (
	"fmt"
	"os"
)

func tearDown(ts *TestSuite) error {
	fmt.Println("tearing down test suite in", ts.TempDir)
	err := os.RemoveAll(ts.TempDir)
	if err != nil {
		return err
	}

	fmt.Println("removing test suite file", testSuiteFile)
	err = os.Remove(testSuiteFile)
	if err != nil {
		return err
	}

	return nil
}