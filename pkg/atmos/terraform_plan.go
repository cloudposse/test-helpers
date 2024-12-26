package atmos

import (
	"fmt"

	"github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/stretchr/testify/require"
)

// Plan runs terraform plan with the given options and returns stdout/stderr.
// This will fail the test if there is an error in the command.
func Plan(t testing.TestingT, options *Options) string {
	out, err := PlanE(t, options)
	require.NoError(t, err)
	return out
}

// PlanE runs terraform plan with the given options and returns stdout/stderr.
func PlanE(t testing.TestingT, options *Options) (string, error) {
	if options.Component == "" {
		return "", ErrorComponentRequired
	}

	if options.Stack == "" {
		return "", ErrorStackRequired
	}
	return RunAtmosCommandE(t, options, FormatArgs(options, "terraform", "plan", "-input=false", "-lock=false")...)
}

// PlanExitCode runs terraform plan with the given options and returns the detailed exitcode.
// This will fail the test if there is an error in the command.
func PlanExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanExitCodeE runs terraform plan with the given options and returns the detailed exitcode.
func PlanExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return GetExitCodeForAtmosCommandE(t, options, FormatArgs(options, "terraform", "plan", "-input=false", "-detailed-exitcode")...)
}

// Custom errors
var (
	ErrorComponentRequired    = fmt.Errorf("you must set ComponentName on options struct to use this function")
	ErrorPlanFilePathRequired = fmt.Errorf("you must set PlanFilePath on options struct to use this function")
	ErrorStackRequired        = fmt.Errorf("you must set StackName on options struct to use this function")
)
