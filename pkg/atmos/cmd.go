package atmos

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	tt "github.com/cloudposse/test-helpers/pkg/testing"
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
	"plan",
	"apply",
	"destroy",
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

	if options.Parallelism > 0 && len(args) > 0 && args[0] == "terraform" && collections.ListContains(commandsWithParallelism, args[1]) {
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
func RunAtmosCommand(t tt.TestingT, additionalOptions *Options, args ...string) string {
	out, err := RunAtmosCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunAtmosCommandE runs atmos with the given arguments and options and return stdout/stderr.
func RunAtmosCommandE(t tt.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.AtmosBinary, args)

	return retry.DoWithRetryableErrorsE(t.(testing.TestingT), description, options.RetryableAtmosErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		s, err := shell.RunCommandAndGetOutputE(t.(testing.TestingT), cmd)
		if err != nil {
			return s, err
		}

		if err := hasWarning(additionalOptions, s); err != nil {
			return s, err
		}
		return s, err
	})

}

// RunAtmosCommandAndGetStdoutE runs atmos with the given arguments and options and returns solely its stdout (but not
// stderr).
func RunAtmosCommandAndGetStdoutE(t tt.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
	return retry.DoWithRetryableErrorsE(t, description, options.RetryableAtmosErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		s, err := shell.RunCommandAndGetStdOutE(t, cmd)
		if err != nil {
			return s, err
		}

		if err := hasWarning(additionalOptions, s); err != nil {
			return s, err
		}

		return s, err
	})
}

// GetExitCodeForAtmosCommand runs atmos with the given arguments and options and returns exit code
func GetExitCodeForAtmosCommand(t tt.TestingT, additionalOptions *Options, args ...string) int {
	exitCode, err := GetExitCodeForAtmosCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// GetExitCodeForAtmosCommandE runs atmos with the given arguments and options and returns exit code
func GetExitCodeForAtmosCommandE(t tt.TestingT, additionalOptions *Options, additionalArgs ...string) (int, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	additionalOptions.Logger.Logf(t, "Running %s with args %v", options.AtmosBinary, options.AtmosBinary, args)
	cmd := generateCommand(options, args...)
	cmd.WorkingDir = options.AtmosBasePath

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

func hasWarning(opts *Options, out string) error {
	for k, v := range opts.WarningsAsErrors {
		str := fmt.Sprintf("\nWarning: %s[^\n]*\n", k)
		re, err := regexp.Compile(str)
		if err != nil {
			return fmt.Errorf("cannot compile regex for warning detection: %w", err)
		}
		m := re.FindAllString(out, -1)
		if len(m) == 0 {
			continue
		}
		return fmt.Errorf("warning(s) were found: %s:\n%s", v, strings.Join(m, ""))
	}
	return nil
}
