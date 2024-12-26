package aws_component_helper

//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/cloudposse/test-helpers/pkg/atmos"
//)
//
//type TestSuite struct {
//	AtmosOptions                  *atmos.Options
//	AwsAccountId                  string
//	AwsRegion                     string
//	ComponentName                 string
//	ComponentSrcPath              string
//	Dependencies                  []*ComponentDeployment
//	FixturesPath                  string
//	ForceNewSuite                 bool
//	Index                         int
//	RandomIdentifier              string
//	SkipSetupComponentUnderTest   bool
//	SkipDeployDependencies        bool
//	SkipDeployComponentUnderTest  bool
//	SkipDestroyComponentUnderTest bool
//	SkipDestroyDependencies       bool
//	SkipTeardownTestSuite         bool
//	SkipTests                     bool
//	SkipVendorDependencies        bool
//	SkipVerifyEnabledFlag         bool
//	SkipNukeTestAccount           bool
//	StackName                     string
//	TempDir                       string
//}
//
//type TestSuites struct {
//	Suites []*TestSuite
//}
//
//// Option type represents a configuration option
//type TestSuiteOption func(*TestSuite)
//
//func (ts *TestSuite) AddDependencies(dependencies []string) {
//	for _, dependency := range dependencies {
//		ts.AddDependency(dependency)
//	}
//}
//
//func (ts *TestSuite) AddCrossStackDependencies(dependencies []ComponentDeployment) {
//	for _, dependency := range dependencies {
//		ts.AddCrossStackDependency(dependency.ComponentName, dependency.StackName)
//	}
//}
//
//func (ts *TestSuite) AddDependency(componentName string) {
//	ts.Dependencies = append(ts.Dependencies, NewAtmosComponent(componentName, ts.StackName))
//}
//
//func (ts *TestSuite) AddCrossStackDependency(componentName string, stackName string) {
//	ts.Dependencies = append(ts.Dependencies, NewAtmosComponent(componentName, stackName))
//}
//
//func (ts *TestSuite) SetupTestSuite(t *testing.T) error {
//	fmt.Println("SetupTestSuite")
//	err := setupTestSuite(ts)
//	return err
//}
//
//func (ts *TestSuite) SetupComponentUnderTest(t *testing.T) error {
//	if !ts.SkipSetupComponentUnderTest {
//		fmt.Println("SetupComponentUnderTest")
//		err := setupComponentUnderTest(ts)
//		return err
//	} else {
//		fmt.Println("Skipping SetupComponentUnderTest")
//	}
//	return nil
//}
//
//func (ts *TestSuite) VendorDependencies(t *testing.T) error {
//	if !ts.SkipVendorDependencies {
//		fmt.Println("VendorDependencies")
//		err := vendorDependencies(t, ts)
//		return err
//	} else {
//		fmt.Println("Skipping VendorDependencies")
//	}
//	return nil
//}
//
//func (ts *TestSuite) DeployDependencies(t *testing.T) error {
//	if !ts.SkipDeployDependencies {
//		fmt.Println("DeployDependencies")
//		err := deployDependencies(t, ts)
//		return err
//	} else {
//		fmt.Println("Skipping DeployDependencies")
//	}
//	return nil
//}
//
//func (ts *TestSuite) VerifyEnabledFlag(t *testing.T) error {
//	if !ts.SkipVerifyEnabledFlag {
//		fmt.Println("VerifyEnabledFlag")
//		_, err := verifyEnabledFlag(t, ts, ts.ComponentName, ts.StackName)
//		return err
//	} else {
//		fmt.Println("Skipping VerifyEnabledFlag")
//	}
//	return nil
//}
//
//func (ts *TestSuite) DeployComponentUnderTest(t *testing.T, vars map[string]interface{}) (string, error) {
//	if !ts.SkipDeployComponentUnderTest {
//		fmt.Println("DeployComponentUnderTest")
//		options, out, err := DeployComponent(t, ts, ts.ComponentName, ts.StackName, vars)
//		ts.AtmosOptions = options
//
//		return out, err
//	} else {
//		fmt.Println("Skipping DeployComponentUnderTest")
//		return "", nil
//	}
//}
//
//func (ts *TestSuite) DestroyComponentUnderTest(t *testing.T, vars map[string]interface{}) (string, error) {
//	if !ts.SkipDestroyComponentUnderTest {
//		fmt.Println("DestroyComponentUnderTest")
//		options, out, err := DestroyComponent(t, ts, ts.ComponentName, ts.StackName, vars)
//		ts.AtmosOptions = options
//
//		return out, err
//	} else {
//		fmt.Println("Skipping DestroyComponentUnderTest")
//		return "", nil
//	}
//}
//
//func (ts *TestSuite) DestroyDependencies(t *testing.T) error {
//	if !skipDestroyDependencies(ts) {
//		fmt.Println("DestroyDependencies")
//		err := destroyDependencies(t, ts)
//		return err
//	} else {
//		fmt.Println("Skipping DestroyDependencies")
//	}
//	return nil
//}
//
//func (ts *TestSuite) TearDownTestSuite(t *testing.T) error {
//	fmt.Println("TeardownTestSuite")
//	if !skipTeardownTestSuite(ts) {
//		err := tearDown(ts)
//		return err
//	} else {
//		fmt.Println("Skipping TeardownTestSuite")
//	}
//	return nil
//}
//
//func (ts *TestSuite) NukeTestAccount(t *testing.T) error {
//	if !skipNukeTestAccount(ts) {
//		fmt.Println("NukeTestAccount")
//		//awsnuke.NukeTestAccountByTag(t, "CreatedByTerratestRun", ts.RandomSeed, []string{ts.AwsRegion}, false)
//	} else {
//		fmt.Println("Skipping NukeTestAccount")
//	}
//
//	return nil
//}
//
//func (ts *TestSuite) Setup(t *testing.T) error {
//	fmt.Println("=== RUN   Test Suite Setup")
//	if err := ts.SetupTestSuite(t); err != nil {
//		return err
//	}
//
//	if err := ts.SetupComponentUnderTest(t); err != nil {
//		return err
//	}
//
//	if err := ts.VendorDependencies(t); err != nil {
//		return err
//	}
//
//	if err := ts.DeployDependencies(t); err != nil {
//		return err
//	}
//
//	if err := ts.VerifyEnabledFlag(t); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (ts *TestSuite) TearDown(t *testing.T) error {
//	fmt.Println("=== RUN   Test Suite TearDown")
//	if err := ts.DestroyDependencies(t); err != nil {
//		return err
//	}
//
//	if err := ts.TearDownTestSuite(t); err != nil {
//		return err
//	}
//
//	if err := ts.NukeTestAccount(t); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func NewTestSuite(awsRegion string, componentName string, stackName string, opts ...TestSuiteOption) (*TestSuite, error) {
//	awsAccountId, err := getAwsAccountId()
//	if err != nil {
//		return &TestSuite{}, err
//	}
//
//	suite := &TestSuite{
//		AtmosOptions:     &atmos.Options{},
//		AwsAccountId:     awsAccountId,
//		AwsRegion:        awsRegion,
//		ComponentName:    componentName,
//		ComponentSrcPath: "../src",
//		FixturesPath:     "fixtures",
//		StackName:        stackName,
//	}
//
//	testName, err := getTestName()
//	if err != nil {
//		return &TestSuite{}, err
//	}
//
//	// Apply optional configurations
//	for _, opt := range opts {
//		opt(suite)
//	}
//
//	// Parse the CLI args
//	suite = parseCLIArgs(suite)
//
//	// Read or create the test suite file
//	suite, err = readOrCreateTestSuiteFile(suite, testName)
//	if err != nil {
//		panic("Failed to create test suite: " + err.Error())
//	}
//
//	return suite, nil
//}
//
//func WithComponentSrcPath(componentSrcPath string) TestSuiteOption {
//	return func(a *TestSuite) {
//		a.ComponentSrcPath = componentSrcPath
//	}
//}
//
//func WithFixturesPath(fixturesPath string) TestSuiteOption {
//	return func(a *TestSuite) {
//		a.FixturesPath = fixturesPath
//	}
//}
//
//func WithDependency(dependency *ComponentDeployment) TestSuiteOption {
//	return func(a *TestSuite) {
//		a.Dependencies = append(a.Dependencies, dependency)
//	}
//}
//
//func WithDependencies(dependencies []*ComponentDeployment) TestSuiteOption {
//	return func(a *TestSuite) {
//		a.Dependencies = append(a.Dependencies, dependencies...)
//	}
//}
