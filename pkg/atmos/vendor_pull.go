package atmos

import (
	"github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/stretchr/testify/require"
)

// Vendor runs atmos vendor with the given options and return stdout/stderr.
func VendorPull(t testing.TestingT, options *Options) string {
	out, err := VendorPullE(t, options)
	require.NoError(t, err)
	return out
}

// VendorPullE runs atmos vendor with the given options and return stdout/stderr.
func VendorPullE(t testing.TestingT, options *Options) (string, error) {
	return RunAtmosCommandE(t, options, FormatArgs(options, "vendor", "pull")...)
}
