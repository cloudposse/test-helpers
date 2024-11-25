package aws_component_helper

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCLIArgs(t *testing.T) {
	// Save original args and restore after test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		initial  *TestSuite
		expected *TestSuite
	}{
		{
			name: "default values not changed when no flags provided",
			args: []string{"prog"},
			initial: &TestSuite{
				SkipNukeTestAccount:           false,
				SkipDeployDependencies:        false,
				SkipDestroyDependencies:       false,
				SkipSetupComponentUnderTest:   false,
				SkipDeployComponentUnderTest:  false,
				SkipDestroyComponentUnderTest: false,
				SkipTeardownTestSuite:         false,
				SkipVerifyEnabledFlag:         false,
			},
			expected: &TestSuite{
				SkipNukeTestAccount:           false,
				SkipDeployDependencies:        false,
				SkipDestroyDependencies:       false,
				SkipSetupComponentUnderTest:   false,
				SkipDeployComponentUnderTest:  false,
				SkipDestroyComponentUnderTest: false,
				SkipTeardownTestSuite:         false,
				SkipVerifyEnabledFlag:         false,
			},
		},
		{
			name: "all flags set to true",
			args: []string{"prog",
				"-skip-aws-nuke",
				"-skip-deploy-deps",
				"-skip-destroy-deps",
				"-skip-setup-cut",
				"-skip-deploy-cut",
				"-skip-destroy-cut",
				"-skip-teardown",
				"-skip-verify-enabled-flag",
			},
			initial: &TestSuite{},
			expected: &TestSuite{
				SkipNukeTestAccount:           true,
				SkipDeployDependencies:        true,
				SkipDestroyDependencies:       true,
				SkipSetupComponentUnderTest:   true,
				SkipDeployComponentUnderTest:  true,
				SkipDestroyComponentUnderTest: true,
				SkipTeardownTestSuite:         true,
				SkipVerifyEnabledFlag:         true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ExitOnError)
			os.Args = tt.args

			result := parseCLIArgs(tt.initial)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSkipDestroyDependencies(t *testing.T) {
	tests := []struct {
		name     string
		ts       *TestSuite
		expected bool
	}{
		{
			name: "both flags false",
			ts: &TestSuite{
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: false,
			},
			expected: false,
		},
		{
			name: "SkipDestroyDependencies true",
			ts: &TestSuite{
				SkipDestroyDependencies:       true,
				SkipDestroyComponentUnderTest: false,
			},
			expected: true,
		},
		{
			name: "SkipDestroyComponentUnderTest true",
			ts: &TestSuite{
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skipDestroyDependencies(tt.ts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSkipTeardownTestSuite(t *testing.T) {
	tests := []struct {
		name     string
		ts       *TestSuite
		expected bool
	}{
		{
			name: "both flags false",
			ts: &TestSuite{
				SkipTeardownTestSuite: false,
			},
			expected: false,
		},
		{
			name: "SkipTeardownTestSuite true",
			ts: &TestSuite{
				SkipTeardownTestSuite: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skipTeardownTestSuite(tt.ts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
