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
				SkipVendorDependencies:        false,
				SkipVerifyEnabledFlag:         false,
				SkipTests:                     false,
				ForceNewSuite:                 false,
				Index:                         -1,
			},
			expected: &TestSuite{
				SkipNukeTestAccount:           false,
				SkipDeployDependencies:        false,
				SkipDestroyDependencies:       false,
				SkipSetupComponentUnderTest:   false,
				SkipDeployComponentUnderTest:  false,
				SkipDestroyComponentUnderTest: false,
				SkipTeardownTestSuite:         false,
				SkipVendorDependencies:        false,
				SkipVerifyEnabledFlag:         false,
				SkipTests:                     false,
				ForceNewSuite:                 false,
				Index:                         -1,
			},
		},
		{
			name: "all skip flags set to true",
			args: []string{"prog",
				"-skip-aws-nuke",
				"-skip-deploy-deps",
				"-skip-destroy-deps",
				"-skip-setup-cut",
				"-skip-deploy-cut",
				"-skip-destroy-cut",
				"-skip-teardown",
				"-skip-vendor",
				"-skip-verify-enabled-flag",
				"-skip-tests",
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
				SkipVendorDependencies:        true,
				SkipVerifyEnabledFlag:         true,
				SkipTests:                     true,
				ForceNewSuite:                 false,
				Index:                         -1,
			},
		},
		{
			name: "force new suite and suite index",
			args: []string{"prog",
				"-force-new-suite",
				"-suite-index=5",
			},
			initial: &TestSuite{},
			expected: &TestSuite{
				ForceNewSuite: true,
				Index:         5,
			},
		},
		{
			name: "initial values respected when not overridden",
			args: []string{"prog",
				"-skip-aws-nuke",
				"-skip-deploy-deps",
			},
			initial: &TestSuite{
				SkipDestroyDependencies:       true,
				SkipSetupComponentUnderTest:   true,
				SkipDeployComponentUnderTest:  true,
				SkipDestroyComponentUnderTest: true,
				SkipTeardownTestSuite:         true,
				SkipVendorDependencies:        true,
				SkipVerifyEnabledFlag:         true,
				SkipTests:                     true,
			},
			expected: &TestSuite{
				SkipNukeTestAccount:           true,
				SkipDeployDependencies:        true,
				SkipDestroyDependencies:       true,
				SkipSetupComponentUnderTest:   true,
				SkipDeployComponentUnderTest:  true,
				SkipDestroyComponentUnderTest: true,
				SkipTeardownTestSuite:         true,
				SkipVendorDependencies:        true,
				SkipVerifyEnabledFlag:         true,
				SkipTests:                     true,
				ForceNewSuite:                 false,
				Index:                         -1,
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
		{
			name: "both flags true",
			ts: &TestSuite{
				SkipDestroyDependencies:       true,
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
			name: "all flags false",
			ts: &TestSuite{
				SkipTeardownTestSuite:         false,
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: false,
			},
			expected: false,
		},
		{
			name: "SkipTeardownTestSuite true",
			ts: &TestSuite{
				SkipTeardownTestSuite:         true,
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: false,
			},
			expected: true,
		},
		{
			name: "SkipDestroyDependencies true",
			ts: &TestSuite{
				SkipTeardownTestSuite:         false,
				SkipDestroyDependencies:       true,
				SkipDestroyComponentUnderTest: false,
			},
			expected: true,
		},
		{
			name: "SkipDestroyComponentUnderTest true",
			ts: &TestSuite{
				SkipTeardownTestSuite:         false,
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: true,
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

func TestSkipNukeTestAccount(t *testing.T) {
	tests := []struct {
		name     string
		ts       *TestSuite
		expected bool
	}{
		{
			name: "all flags false",
			ts: &TestSuite{
				SkipNukeTestAccount:           false,
				SkipTeardownTestSuite:         false,
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: false,
			},
			expected: false,
		},
		{
			name: "SkipNukeTestAccount true",
			ts: &TestSuite{
				SkipNukeTestAccount:           true,
				SkipTeardownTestSuite:         false,
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: false,
			},
			expected: true,
		},
		{
			name: "SkipTeardownTestSuite true",
			ts: &TestSuite{
				SkipNukeTestAccount:           false,
				SkipTeardownTestSuite:         true,
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: false,
			},
			expected: true,
		},
		{
			name: "SkipDestroyDependencies true",
			ts: &TestSuite{
				SkipNukeTestAccount:           false,
				SkipTeardownTestSuite:         false,
				SkipDestroyDependencies:       true,
				SkipDestroyComponentUnderTest: false,
			},
			expected: true,
		},
		{
			name: "SkipDestroyComponentUnderTest true",
			ts: &TestSuite{
				SkipNukeTestAccount:           false,
				SkipTeardownTestSuite:         false,
				SkipDestroyDependencies:       false,
				SkipDestroyComponentUnderTest: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skipNukeTestAccount(tt.ts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
