package atmos

import (
	"fmt"
	"os/exec"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

func generateCommand(options *Options, args ...string) shell.Command {
	cmd := shell.Command{
		Command:    options.AtmosBinary,
		Args:       args,
		WorkingDir: options.AtmosBasePath,
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}
	return cmd
}

var commandsWithParallelism = []string{
	"terraform plan",
	"terraform apply",
	"terraform destroy",
}

const (
	// AtmosDefaultPath command to run atmos
	AtmosDefaultPath = "atmos"
)

var DefaultExecutable = defaultAtmosExecutable()

// GetCommonOptions extracts commons atmos options
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	if options.AtmosBinary == "" {
		options.AtmosBinary = DefaultExecutable
	}

	if options.Parallelism > 0 && len(args) > 0 && collections.ListContains(commandsWithParallelism, args[0]) {
		args = append(args, fmt.Sprintf("--parallelism=%d", options.Parallelism))
	}

	// if SshAgent is provided, override the local SSH agent with the socket of our in-process agent
	if options.SshAgent != nil {
		// Initialize EnvVars, if it hasn't been set yet
		if options.EnvVars == nil {
			options.EnvVars = map[string]string{}
		}
		options.EnvVars["SSH_AUTH_SOCK"] = options.SshAgent.SocketFile()
	}
	return options, args
}

// RunAtmosCommand runs atmos with the given arguments and options and return stdout/stderr.
func RunAtmosCommand(t testing.TestingT, additionalOptions *Options, args ...string) string {
	out, err := RunAtmosCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunAtmosCommandE runs atmos with the given arguments and options and return stdout/stderr.
func RunAtmosCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
	return retry.DoWithRetryableErrorsE(t, description, options.RetryableAtmosErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandAndGetOutputE(t, cmd)
	})
}

// RunAtmosCommandAndGetStdoutE runs atmos with the given arguments and options and returns solely its stdout (but not
// stderr).
func RunAtmosCommandAndGetStdoutE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
	return retry.DoWithRetryableErrorsE(t, description, options.RetryableAtmosErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandAndGetStdOutE(t, cmd)
	})
}

// GetExitCodeForAtmosCommand runs atmos with the given arguments and options and returns exit code
func GetExitCodeForAtmosCommand(t testing.TestingT, additionalOptions *Options, args ...string) int {
	exitCode, err := GetExitCodeForAtmosCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// GetExitCodeForAtmosCommandE runs atmos with the given arguments and options and returns exit code
func GetExitCodeForAtmosCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (int, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	additionalOptions.Logger.Logf(t, "Running %s with args %v", options.AtmosBinary, args)
	cmd := generateCommand(options, args...)
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err == nil {
		return DefaultSuccessExitCode, nil
	}
	exitCode, getExitCodeErr := shell.GetExitCodeForRunCommandError(err)
	if getExitCodeErr == nil {
		return exitCode, nil
	}
	return DefaultErrorExitCode, getExitCodeErr
}

func defaultAtmosExecutable() string {
	cmd := exec.Command(AtmosDefaultPath, "-version")
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err == nil {
		return AtmosDefaultPath
	}

	return AtmosDefaultPath
}
