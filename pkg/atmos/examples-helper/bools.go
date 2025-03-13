package examples_helper

import "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"

func anyPhasesSkipped(config *config.Config) bool {
	return config.OnlyDeployDependencies ||
		config.SkipDeployComponent ||
		config.SkipDeployDependencies ||
		config.SkipDestroyComponent ||
		config.SkipDestroyDependencies ||
		config.SkipSetupTestSuite ||
		config.SkipTeardownTestSuite ||
		config.SkipVendorDependencies
}
