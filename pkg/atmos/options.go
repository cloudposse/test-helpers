package atmos

import (
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
)

var (
	DefaultRetryableAtmosErrors = map[string]string{
		// Helm related atmos calls may fail when too many tests run in parallel. While the exact cause is unknown,
		// this is presumably due to all the network contention involved. Usually a retry resolves the issue.
		".*read: connection reset by peer.*": "Failed to reach helm charts repository.",
		".*transport is closing.*":           "Failed to reach Kubernetes API.",

		// `atmos init` frequently fails in CI due to network issues accessing plugins. The reason is unknown, but
		// eventually these succeed after a few retries.
		".*unable to verify signature.*":                  "Failed to retrieve plugin due to transient network error.",
		".*unable to verify checksum.*":                   "Failed to retrieve plugin due to transient network error.",
		".*no provider exists with the given name.*":      "Failed to retrieve plugin due to transient network error.",
		".*registry service is unreachable.*":             "Failed to retrieve plugin due to transient network error.",
		".*Error installing provider.*":                   "Failed to retrieve plugin due to transient network error.",
		".*Failed to query available provider packages.*": "Failed to retrieve plugin due to transient network error.",
		".*timeout while waiting for plugin to start.*":   "Failed to retrieve plugin due to transient network error.",
		".*timed out waiting for server handshake.*":      "Failed to retrieve plugin due to transient network error.",
		"could not query provider registry for":           "Failed to retrieve plugin due to transient network error.",

		// Provider bugs where the data after apply is not propagated. This is usually an eventual consistency issue, so
		// retrying should self resolve it.
		// See https://github.com/atmos-providers/atmos-provider-aws/issues/12449 for an example.
		".*Provider produced inconsistent result after apply.*": "Provider eventual consistency error.",
	}
)

// Options for running Atmos commands
type Options struct {
	AtmosBinary   string // Name of the binary that will be used
	AtmosBasePath string // The path of the atmos root for components and stacks

	// The vars to pass to Atmos commands using the -var option. Note that atmos does not support passing `null`
	// as a variable value through the command line. That is, if you use `map[string]interface{}{"foo": nil}` as `Vars`,
	// this will translate to the string literal `"null"` being assigned to the variable `foo`. However, nulls in
	// lists and maps/objects are supported. E.g., the following var will be set as expected (`{ bar = null }`:
	// map[string]interface{}{
	//     "foo": map[string]interface{}{"bar": nil},
	// }
	Vars                      map[string]interface{}
	VarFiles                  []string               // The var file paths to pass to Atmos commands using -var-file option.
	Stack                     string                 // The stack to pass to the atmos command with -stack
	Component                 string                 // The component to pass to the atmos terraform command
	Targets                   []string               // The target resources to pass to the atmos command with -target
	Lock                      bool                   // The lock option to pass to the atmos command with -lock
	LockTimeout               string                 // The lock timeout option to pass to the atmos command with -lock-timeout
	EnvVars                   map[string]string      // Environment variables to set when running Atmos
	BackendConfig             map[string]interface{} // The vars to pass to the atmos init command for extra configuration for the backend
	RetryableAtmosErrors      map[string]string      // If Atmos apply fails with one of these (transient) errors, retry. The keys are a regexp to match against the error and the message is what to display to a user if that error is matched.
	MaxRetries                int                    // Maximum number of times to retry errors matching RetryableAtmosErrors
	TimeBetweenRetries        time.Duration          // The amount of time to wait between retries
	Upgrade                   bool                   // Whether the -upgrade flag of the atmos init command should be set to true or not
	Reconfigure               bool                   // Set the -reconfigure flag to the atmos init command
	MigrateState              bool                   // Set the -migrate-state and -force-copy (suppress 'yes' answer prompt) flag to the atmos init command
	NoColor                   bool                   // Whether the -no-color flag will be set for any Atmos command or not
	RedirectStrErrDestination string                 // The destination to redirect stderr to. If nil, stderr will not be redirected.
	SshAgent                  *ssh.SshAgent          // Overrides local SSH agent with the given in-process agent
	NoStderr                  bool                   // Disable stderr redirection
	OutputMaxLineSize         int                    // The max size of one line in stdout and stderr (in bytes)
	Logger                    *logger.Logger         // Set a non-default logger that should be used. See the logger package for more info.
	Parallelism               int                    // Set the parallelism setting for Atmos
	PlanFilePath              string                 // The path to output a plan file to (for the plan command) or read one from (for the apply command)
	PluginDir                 string                 // The path of downloaded plugins to pass to the atmos init command (-plugin-dir)
	SetVarsAfterVarFiles      bool                   // Pass -var options after -var-file options to Atmos commands
}

// Clone makes a deep copy of most fields on the Options object and returns it.
//
// NOTE: options.SshAgent and options.Logger CANNOT be deep copied (e.g., the SshAgent struct contains channels and
// listeners that can't be meaningfully copied), so the original values are retained.
func (options *Options) Clone() (*Options, error) {
	newOptions := &Options{}
	if err := copier.Copy(newOptions, options); err != nil {
		return nil, err
	}
	// copier does not deep copy maps, so we have to do it manually.
	newOptions.EnvVars = make(map[string]string)
	for key, val := range options.EnvVars {
		newOptions.EnvVars[key] = val
	}
	newOptions.Vars = make(map[string]interface{})
	for key, val := range options.Vars {
		newOptions.Vars[key] = val
	}
	newOptions.BackendConfig = make(map[string]interface{})
	for key, val := range options.BackendConfig {
		newOptions.BackendConfig[key] = val
	}
	newOptions.RetryableAtmosErrors = make(map[string]string)
	for key, val := range options.RetryableAtmosErrors {
		newOptions.RetryableAtmosErrors[key] = val
	}

	return newOptions, nil
}

// WithDefaultRetryableErrors makes a copy of the Options object and returns an updated object with sensible defaults
// for retryable errors. The included retryable errors are typical errors that most atmos modules encounter during
// testing, and are known to self resolve upon retrying. This will fail the test if there are any errors in the cloning
// process.
func WithDefaultRetryableErrors(t testing.TestingT, originalOptions *Options) *Options {
	newOptions, err := originalOptions.Clone()
	require.NoError(t, err)

	if newOptions.RetryableAtmosErrors == nil {
		newOptions.RetryableAtmosErrors = map[string]string{}
	}
	for k, v := range DefaultRetryableAtmosErrors {
		newOptions.RetryableAtmosErrors[k] = v
	}

	newOptions.MaxRetries = 3
	newOptions.TimeBetweenRetries = 5 * time.Second

	return newOptions
}
