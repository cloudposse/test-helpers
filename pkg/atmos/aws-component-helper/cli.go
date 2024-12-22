package aws_component_helper

//func parseCLIArgs(ts *TestSuite) *TestSuite {
//	forceNewSuite := flag.Bool("force-new-suite", false, "force new suite")
//	suiteIndex := flag.Int("suite-index", -1, "suite index")
//	skipAwsNuke := flag.Bool("skip-aws-nuke", ts.SkipNukeTestAccount, "skip aws nuke")
//	skipDeployDependencies := flag.Bool("skip-deploy-deps", ts.SkipDeployDependencies, "skip deploy dependencies")
//	skipDestroyDependencies := flag.Bool("skip-destroy-deps", ts.SkipDestroyDependencies, "skip destroy dependencies")
//	skipSetupComponentUnderTest := flag.Bool("skip-setup-cut", ts.SkipSetupComponentUnderTest, "skip setup component under test")
//	skipDeployComponentUnderTest := flag.Bool("skip-deploy-cut", ts.SkipDeployComponentUnderTest, "skip deploy component under test")
//	skipDestroyComponentUnderTest := flag.Bool("skip-destroy-cut", ts.SkipDestroyComponentUnderTest, "skip destroy component under test")
//	skipTeardownTestSuite := flag.Bool("skip-teardown", ts.SkipTeardownTestSuite, "skip test suite teardown")
//	skipVendorDependencies := flag.Bool("skip-vendor", ts.SkipVendorDependencies, "skip vendor dependencies")
//	skipVerifyEnabledFlag := flag.Bool("skip-verify-enabled-flag", ts.SkipVerifyEnabledFlag, "skip verify enabled flag")
//	skipTests := flag.Bool("skip-tests", ts.SkipTests, "skip tests")
//
//	flag.Parse()
//
//	ts.ForceNewSuite = *forceNewSuite
//	ts.Index = *suiteIndex
//	ts.SkipNukeTestAccount = *skipAwsNuke
//	ts.SkipDeployDependencies = *skipDeployDependencies
//	ts.SkipDestroyDependencies = *skipDestroyDependencies
//	ts.SkipSetupComponentUnderTest = *skipSetupComponentUnderTest
//	ts.SkipDeployComponentUnderTest = *skipDeployComponentUnderTest
//	ts.SkipDestroyComponentUnderTest = *skipDestroyComponentUnderTest
//	ts.SkipTeardownTestSuite = *skipTeardownTestSuite
//	ts.SkipVendorDependencies = *skipVendorDependencies
//	ts.SkipVerifyEnabledFlag = *skipVerifyEnabledFlag
//	ts.SkipTests = *skipTests
//	return ts
//}
//
//func skipDestroyDependencies(ts *TestSuite) bool {
//	return ts.SkipDestroyDependencies || ts.SkipDestroyComponentUnderTest
//}
//
//func skipTeardownTestSuite(ts *TestSuite) bool {
//	return ts.SkipTeardownTestSuite || skipDestroyDependencies(ts)
//}
//
//func skipNukeTestAccount(ts *TestSuite) bool {
//	return ts.SkipNukeTestAccount || skipTeardownTestSuite(ts)
//}
