package aws_component_helper

//func tearDown(ts *TestSuite) error {
//	fmt.Println("tearing down test suite in", ts.TempDir)
//	err := os.RemoveAll(ts.TempDir)
//	if err != nil {
//		return err
//	}
//
//	fmt.Println("removing test suite file", testSuiteFile)
//	err = os.Remove(testSuiteFile)
//	if err != nil {
//		return err
//	}
//
//	defer os.Unsetenv("ATMOS_BASE_PATH")
//	defer os.Unsetenv("ATMOS_CLI_CONFIG_PATH")
//	defer os.Unsetenv("TEST_ACCOUNT_ID")
//
//	return nil
//}
