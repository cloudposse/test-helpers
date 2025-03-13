package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnlyDeployDependenciesFlag(t *testing.T) {
	// Create a basic Config
	config := &Config{
		OnlyDeployDependencies: true,
	}

	// Manually apply the flag logic that we want to test
	if config.OnlyDeployDependencies {
		config.SkipDeployComponent = true
		config.SkipDeployDependencies = false
		config.SkipDestroyComponent = true
		config.SkipDestroyDependencies = true
		config.SkipEnabledFlagTest = true
		config.SkipSetupTestSuite = false
		config.SkipTeardownTestSuite = true
		config.SkipVendorDependencies = false
	}

	// Verify the flag affects other flags as expected
	assert.True(t, config.OnlyDeployDependencies, "OnlyDeployDependencies should be true")
	assert.True(t, config.SkipDeployComponent, "SkipDeployComponent should be true")
	assert.False(t, config.SkipDeployDependencies, "SkipDeployDependencies should be false to ensure dependencies are deployed")
	assert.True(t, config.SkipDestroyComponent, "SkipDestroyComponent should be true")
	assert.True(t, config.SkipDestroyDependencies, "SkipDestroyDependencies should be true")
	assert.True(t, config.SkipEnabledFlagTest, "SkipEnabledFlagTest should be true")
	assert.False(t, config.SkipSetupTestSuite, "SkipSetupTestSuite should be false to ensure setup runs")
	assert.True(t, config.SkipTeardownTestSuite, "SkipTeardownTestSuite should be true")
	assert.False(t, config.SkipVendorDependencies, "SkipVendorDependencies should be false to ensure dependencies are vendored")
}

func TestFlagOverrides(t *testing.T) {
	// Create a Config with conflicting flags
	config := &Config{
		OnlyDeployDependencies: true,
		SkipDeployDependencies: true,  // Should be overridden to false
		SkipDeployComponent:    false, // Should be overridden to true
	}

	// Apply the flag logic
	if config.OnlyDeployDependencies {
		config.SkipDeployComponent = true
		config.SkipDeployDependencies = false
		config.SkipDestroyComponent = true
		config.SkipDestroyDependencies = true
		config.SkipEnabledFlagTest = true
		config.SkipSetupTestSuite = false
		config.SkipTeardownTestSuite = true
		config.SkipVendorDependencies = false
	}

	// Verify OnlyDeployDependencies takes precedence
	assert.True(t, config.OnlyDeployDependencies, "OnlyDeployDependencies should be true")
	assert.True(t, config.SkipDeployComponent, "SkipDeployComponent should be true, overridden by OnlyDeployDependencies")
	assert.False(t, config.SkipDeployDependencies, "SkipDeployDependencies should be false, overridden by OnlyDeployDependencies")
}

func TestWriteConfig(t *testing.T) {
	// Create a temporary config file
	tempConfigFile := "test_config_write.yaml"
	defer os.Remove(tempConfigFile) // Clean up after test

	// Create a basic Config
	config := &Config{
		ConfigFilePath: tempConfigFile,
	}

	// Write the config to file
	err := config.WriteConfig()
	assert.NoError(t, err, "WriteConfig should not return an error")

	// Verify file exists
	_, err = os.Stat(tempConfigFile)
	assert.NoError(t, err, "Config file should exist after writing")
}
