package component_helper

import "github.com/cloudposse/test-helpers/pkg/atmos/component-helper/config"

func anyPhasesSkipped(config *config.Config) bool {
	return config.SkipDeployComponent ||
		config.SkipDeployDependencies ||
		config.SkipDestroyComponent ||
		config.SkipDestroyDependencies ||
		config.SkipSetupTestSuite ||
		config.SkipTeardownTestSuite ||
		config.SkipVendorDependencies
}
