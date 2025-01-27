package config

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func init() {
	flag.String("config", "test_suite.yaml", "The path to the config file")
	flag.String("fixtures-dir", "fixtures", "The path to the fixtures directory")
	flag.Bool("skip-deploy-component", true, "Disables running the deploy component phase of tests")
	flag.Bool("skip-deploy-dependencies", true, "Disables running the deploy dependencies phase of tests")
	flag.Bool("skip-destroy-component", true, "Disables running the destroy component phase of tests")
	flag.Bool("skip-destroy-dependencies", true, "Disables running the destroy dependencies phase of tests")
	flag.Bool("skip-enabled-flag-test", true, "Disables running the Enabled flag test")
	flag.Bool("skip-setup", true, "Disables running the setup test suite phase of tests")
	flag.Bool("skip-teardown", true, "Disables running the teardown test suite phase of tests")
	flag.Bool("skip-vendor", true, "Disables running the vendor dependencies phase of tests")
	flag.String("src-dir", "", "The path to the component source directory")
	flag.String("state-dir", "", "The path to the terraform state directory")
	flag.String("temp-dir", "", "The path to the temp directory")
}

type Config struct {
	ConfigFilePath          string
	FixturesDir             string
	RandomIdentifier        string
	SkipDeployComponent     bool
	SkipDeployDependencies  bool
	SkipDestroyComponent    bool
	SkipDestroyDependencies bool
	SkipEnabledFlagTest     bool
	SkipSetupTestSuite      bool
	SkipTeardownTestSuite   bool
	SkipVendorDependencies  bool
	SrcDir                  string
	StateDir                string
	TempDir                 string
}

func (c *Config) WriteConfig() error {
	return writeConfigWithoutPFlags(c.ConfigFilePath)
}

func writeConfigWithoutPFlags(filename string) error {
	// Create a temporary viper instance
	tempViper := viper.New()

	// Get all settings
	allSettings := viper.AllSettings()

	// Copy settings to tempViper, skipping pflag-bound keys
	for key, value := range allSettings {
		if !isPFlagBound(key) {
			tempViper.Set(key, value)
		}
	}

	// Write the filtered configuration to the file
	return tempViper.WriteConfigAs(filename)
}

// isPFlagBound checks if a key is bound to a pflag
func isPFlagBound(key string) bool {
	if key == "configfilepath" {
		return false
	}

	// Check if the key corresponds to a defined flag
	return strings.HasPrefix(key, "skip") ||
		strings.HasPrefix(key, "test") ||
		flag.Lookup(key) != nil
}

func InitConfig(t *testing.T) *Config {
	viper.SetConfigName("test_suite")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("ConfigFilePath", "test_suite.yaml")
	viper.SetDefault("FixturesDir", "fixtures")

	randID := random.UniqueId()
	viper.SetDefault("RandomIdentifier", strings.ToLower(randID))
	viper.SetDefault("SkipDeployComponent", false)
	viper.SetDefault("SkipDeployDependencies", false)
	viper.SetDefault("SkipDestroyComponent", false)
	viper.SetDefault("SkipDestroyDependencies", false)
	viper.SetDefault("SkipEnabledFlagTest", false)
	viper.SetDefault("SkipSetupTestSuite", false)
	viper.SetDefault("SkipTeardownTestSuite", false)
	viper.SetDefault("SkipVendorDependencies", false)
	viper.SetDefault("TempDir", "")
	viper.SetDefault("SrcDir", "../src")
	viper.SetDefault("StateDir", "")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	require.NoError(t, err)

	err = viper.BindPFlag("ConfigFilePath", pflag.Lookup("config"))
	require.NoError(t, err)

	err = viper.BindPFlag("FixturesDir", pflag.Lookup("fixtures-dir"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipEnabledFlagTest", pflag.Lookup("skip-enabled-flag-test"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipDeployComponent", pflag.Lookup("skip-deploy-component"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipDeployDependencies", pflag.Lookup("skip-deploy-dependencies"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipDestroyComponent", pflag.Lookup("skip-destroy-component"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipDestroyDependencies", pflag.Lookup("skip-destroy-dependencies"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipSetupTestSuite", pflag.Lookup("skip-setup"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipTeardownTestSuite", pflag.Lookup("skip-teardown"))
	require.NoError(t, err)

	err = viper.BindPFlag("StateDir", pflag.Lookup("state-dir"))
	require.NoError(t, err)

	err = viper.BindPFlag("SrcDir", pflag.Lookup("src-dir"))
	require.NoError(t, err)

	err = viper.BindPFlag("TempDir", pflag.Lookup("temp-dir"))
	require.NoError(t, err)

	err = viper.BindPFlag("SkipVendorDependencies", pflag.Lookup("skip-vendor"))
	require.NoError(t, err)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error and use defaults
		} else {
			t.Fatal(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		t.Fatal(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	err = writeConfigWithoutPFlags(viper.GetString("ConfigFilePath"))
	require.NoError(t, err)

	return config
}
