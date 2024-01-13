package atmos

import (
	"errors"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Apply runs atmos terraform apply with the given options and return stdout/stderr. Note that this method does NOT call
// destroy and assumes the caller is responsible for cleaning up any resources created by running apply.
func Apply(t testing.TestingT, options *Options) string {
	out, err := ApplyE(t, options)
	require.NoError(t, err)
	return out
}

// ApplyE runs atmos terraform apply with the given options and return stdout/stderr. Note that this method does NOT
// call destroy and assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyE(t testing.TestingT, options *Options) (string, error) {
	if options.Component == "" {
		return "", ErrorComponentRequired
	}

	if options.Stack == "" {
		return "", ErrorStackRequired
	}

	return RunAtmosCommandE(t, options, FormatArgs(options, "terraform", "apply", "-input=false", "-auto-approve")...)
}

// ApplyAndIdempotent runs atmos terraform apply with the given options and return stdout/stderr from the apply command.
// It then runs plan again and will fail the test if plan requires additional changes. Note that this method does NOT
// call destroy and assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyAndIdempotent(t testing.TestingT, options *Options) string {
	out, err := ApplyAndIdempotentE(t, options)
	require.NoError(t, err)

	return out
}

// ApplyAndIdempotentE runs atmos terraform apply with the given options and return stdout/stderr from the apply
// command. It then runs plan again and will fail the test if plan requires additional changes. Note that this method
// does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyAndIdempotentE(t testing.TestingT, options *Options) (string, error) {
	out, err := ApplyE(t, options)

	if err != nil {
		return out, err
	}

	exitCode, err := PlanExitCodeE(t, options)

	if err != nil {
		return out, err
	}

	if exitCode != 0 {
		return out, errors.New("terraform configuration not idempotent")
	}

	return out, nil
}
