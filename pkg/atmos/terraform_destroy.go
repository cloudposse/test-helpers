package atmos

import (
	"github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/stretchr/testify/require"
)

// Destroy runs atmos terraform destroy with the given options and return stdout/stderr.
func Destroy(t testing.TestingT, options *Options) string {
	out, err := DestroyE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyE runs atmos terraform destroy with the given options and return stdout/stderr.
func DestroyE(t testing.TestingT, options *Options) (string, error) {
	if options.Component == "" {
		return "", ErrorComponentRequired
	}

	if options.Stack == "" {
		return "", ErrorStackRequired
	}

	return RunAtmosCommandE(t, options, FormatArgs(options, "terraform", "destroy", "-input=false", "-auto-approve")...)
}
