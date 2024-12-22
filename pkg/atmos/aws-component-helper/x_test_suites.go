package aws_component_helper

import (
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var (
	getAwsAaccountIdCallback = getAwsAccountId
)

type XTestSuites struct {
	RandomIdentifier string
	AtmosOptions     *atmos.Options
	AwsAccountId     string
	AwsRegion        string
	SourceDir        string
	FixturesPath     string
	WorkDir          string
	suites           map[string]*XTestSuite
}

func NewTestSuites(t *testing.T, sourceDir string, awsRegion string, fixturesDir string) (*XTestSuites, error) {
	awsAccountId, err := getAwsAaccountIdCallback()
	if err != nil {
		return &XTestSuites{}, err
	}

	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	workdir := filepath.Join(os.TempDir(), "test-suites-"+randomId)
	suitesOptions := &atmos.Options{
		AtmosBasePath: filepath.Join(workdir, fixturesDir),
		NoColor:       true,
		Vars: map[string]interface{}{
			"region": awsRegion,
		},
	}

	suites := &XTestSuites{
		RandomIdentifier: randomId,
		SourceDir:        sourceDir,
		FixturesPath:     fixturesDir,
		WorkDir:          workdir,
		AwsAccountId:     awsAccountId,
		AwsRegion:        awsRegion,
		AtmosOptions:     suitesOptions,
		suites:           map[string]*XTestSuite{},
	}

	describeStacksConfigs, err := atmos.DescribeStacksE(t, &atmos.Options{
		AtmosBasePath: filepath.Join(suites.SourceDir, fixturesDir),
		Vars:          map[string]interface{}{},
	})

	if err != nil {
		return suites, err
	}

	for stackName, stack := range describeStacksConfigs.Stacks {
		for componentName, component := range stack.Components.Terraform {
			if component.Settings.Test != nil {
				if component.Settings.Test.Suite == "" {
					return &XTestSuites{}, fmt.Errorf("settings.test.suite is required for component %s in stack %s", componentName, stackName)
				}

				suiteName := component.Settings.Test.Suite

				validStages := []string{"testSuiteSetUp", "testSetUp", "subjectUnderTest", "assert"}

				if component.Settings.Test.Stage == "" {
					return &XTestSuites{}, fmt.Errorf("settings.test.stage is required for component %s in stack %s", componentName, stackName)
				} else if !slices.Contains(validStages, component.Settings.Test.Stage) {
					message := "settings.test.stage should be one of \"testSuiteSetUp\", \"testSetUp\", \"subjectUnderTest\", \"assert\" for component %s in stack %s"
					return &XTestSuites{}, fmt.Errorf(message, componentName, stackName)
				}

				stage := component.Settings.Test.Stage

				if stage != "testSuiteSetUp" && component.Settings.Test.Test == "" {
					return &XTestSuites{}, fmt.Errorf("settings.test.test is required for component %s in stack %s", componentName, stackName)
				}

				testName := component.Settings.Test.Test

				suite := suites.GetOrCreateSuite(suiteName)

				switch stage {
				case "testSuiteSetUp":
					suite.AddSetup(componentName, stackName)
				case "testSetUp":
					suite.GetOrCreateTest(testName).AddSetup(componentName, stackName)
				case "subjectUnderTest":
					suite.GetOrCreateTest(testName).SetSubject(componentName, stackName)
				case "assert":
					suite.GetOrCreateTest(testName).AddSAssert(componentName, stackName)
				}
			}
		}
	}

	//for _, opt := range opts {
	//	opt(suite)
	//}

	return suites, nil
}

func (ts *XTestSuites) Run(t *testing.T) {
	fmt.Printf("create TMP dir: %s \n", ts.WorkDir)

	err := os.Mkdir(ts.WorkDir, 0777)
	assert.NoError(t, err)
	defer os.RemoveAll(ts.WorkDir)

	err = copyDirectoryRecursively(ts.SourceDir, ts.WorkDir)
	assert.NoError(t, err)

	atmos.VendorPull(t, ts.AtmosOptions)

	err = createStateDir(ts.WorkDir)
	assert.NoError(t, err)

	// t.Parallel()
	for name, suite := range ts.suites {
		t.Run(name, func(t *testing.T) {
			suite.Run(t)
		})
	}
}

func (ts *XTestSuites) GetOrCreateSuite(name string) *XTestSuite {
	if _, ok := ts.suites[name]; !ok {
		ts.suites[name] = NewXTestSuite(ts.AtmosOptions)
	}
	return ts.suites[name]

}

func (ts *XTestSuites) SetSuite(name string, item *XTestSuite) {
	ts.suites[name] = item
}
