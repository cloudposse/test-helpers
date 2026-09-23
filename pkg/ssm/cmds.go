// Package ssm supplies a focused set of helper assertions for Terratest
// integration suites.  Each helper runs a shell command on a target EC2
// instance through AWS Systems Manager (SSM) and immediately translates the
// result into a testify assertion, providing declarative checks of
// common OS-level expectations:
//
//   - file or directory presence and permissions
//   - package installation status
//   - user, group, or service existence and state
//   - mount points, network reachability, DNS resolution
//   - arbitrary command output
//
// The helpers hide repetitive SSM polling and error handling, so test code
// stays short, readable, and intent-focused.
package ssm

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
)

// FileExists checks if a file exists at the given path.
func FileExists(t *testing.T, ssmClient *ssm.Client, instanceID string, path string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("test -f %s && printf exists", path), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "exists", "Expected file to exist: "+path)
}

// FileExecutable checks if a file at the given path is executable.
func FileExecutable(t *testing.T, ssmClient *ssm.Client, instanceID string, path string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("test -x %s && printf executable", path), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check file executability: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "executable", "Expected file to be executable: "+path)
}

// DirectoryExists checks if a directory exists at the given path.
func DirectoryExists(t *testing.T, ssmClient *ssm.Client, instanceID string, path string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("test -d %s && printf exists", path), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check directory existence: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "exists", "Expected directory to exist: "+path)
}

// PathIsMount checks if a path is a mount point.
func PathIsMount(t *testing.T, ssmClient *ssm.Client, instanceID string, path string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("mountpoint -q %s && printf mount", path), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check mount point: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "mount", "Expected path to be a mount point: "+path)
}

// PackageInstalled checks if a package is installed.
func PackageInstalled(t *testing.T, ssmClient *ssm.Client, instanceID string, pkg string) {
	t.Helper()
	cmd := fmt.Sprintf(
		`(rpm -q %s >/dev/null 2>&1 && exit 0) || (dpkg-query -W -f='${Status}' %s 2>/dev/null | grep -qx 'install ok installed' && exit 0) || exit 1`,
		pkg, pkg,
	)
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, cmd, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check package installation: %v", err)
	}
	_ = result // non-zero exit code means the package is not installed (already asserted by CheckSSMCommandWithClientE)
}

// UserExists checks if a user exists.
func UserExists(t *testing.T, ssmClient *ssm.Client, instanceID string, username string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("id -nu %s", username), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check user existence: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), username, "Expected user to exist: "+username)
}

// GroupExists checks if a group exists.
func GroupExists(t *testing.T, ssmClient *ssm.Client, instanceID string, groupname string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("getent group %s | awk -F: '{print $1}'", groupname), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check group existence: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), groupname, "Expected group to exist: "+groupname)
}

// ServiceEnabled checks if a systemd service is enabled.
func ServiceEnabled(t *testing.T, ssmClient *ssm.Client, instanceID string, service string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("systemctl is-enabled %s", service), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check service enabled: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "enabled", "Expected service to be enabled: "+service)
}

// ServiceRunning checks if a systemd service is running.
func ServiceRunning(t *testing.T, ssmClient *ssm.Client, instanceID string, service string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("systemctl is-active %s", service), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check service running: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "active", "Expected service to be running: "+service)
}

// ServiceDisabled checks if a systemd service is disabled. The `|| true`
// guard is required because `systemctl is-enabled` exits non-zero for a
// disabled unit, which would otherwise be treated as a failed SSM command.
func ServiceDisabled(t *testing.T, ssmClient *ssm.Client, instanceID string, service string) {
	t.Helper()
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, fmt.Sprintf("systemctl is-enabled %s || true", service), 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check service disabled: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "disabled", "Expected service to be disabled: "+service)
}

// FileContains asserts that a file contains a specific string.
func FileContains(t *testing.T, ssmClient *ssm.Client, instanceID string, path string, substr string) {
	cmd := fmt.Sprintf("cat %s", path)
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, cmd, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check file containment: %v", err)
	}
	assert.Contains(t, result.Stdout, substr, "Expected file %s to contain %q", path, substr)
}

// FileContainsAs asserts that a file contains a specific string, as another user.
func FileContainsAs(t *testing.T, ssmClient *ssm.Client, instanceID string, user string, path string, substr string) {
	cmd := fmt.Sprintf("sudo -u %s -- cat %s", user, path)
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, cmd, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check file containment: %v", err)
	}
	assert.Contains(t, result.Stdout, substr, "Expected file %s to contain %q", path, substr)
}

// CommandOutputContains asserts that a command output contains a specific string.
func CommandOutputContains(t *testing.T, ssmClient *ssm.Client, instanceID string, command string, expected string) {
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, command, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check command: %v", err)
	}
	assert.Contains(t, result.Stdout, expected, "Expected command (%s) to contain %q", command, expected)
}

// RunCommand runs a command and returns the output.
func RunCommand(t *testing.T, ssmClient *ssm.Client, instanceID string, command string) string {
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, command, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check command: %v", err)
	}
	return result.Stdout
}

// RunCommandAs runs a command as another user and returns the output.
func RunCommandAs(t *testing.T, ssmClient *ssm.Client, instanceID string, user string, command string) string {
	wrapped := fmt.Sprintf("sudo -iu %s -- bash -c %q", user, command)
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, wrapped, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check command as user %s: %v", user, err)
	}
	return result.Stdout
}

// PortReachable checks if a host and port are reachable from the EC2 instance.
func PortReachable(t *testing.T, ssmClient *ssm.Client, instanceID string, host string, port int) {
	cmd := fmt.Sprintf("nc -z -w3 %s %d && printf reachable || printf unreachable", host, port)
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, cmd, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check port reachability: %v", err)
	}
	assert.Equal(t, strings.TrimSpace(result.Stdout), "reachable", "Expected %s:%d to be reachable", host, port)
}

// DNSResolves asserts that a hostname resolves to an IP.
func DNSResolves(t *testing.T, ssmClient *ssm.Client, instanceID string, host string, expectedIP string) {
	cmd := fmt.Sprintf("getent hosts %s | awk '{print $1}'", host) // this may need to change to `getent ahosts` because of ipv6
	result, err := aws.CheckSSMCommandWithClientE(t, ssmClient, instanceID, cmd, 1*time.Minute)
	if err != nil {
		t.Fatalf("Failed to check DNS resolution: %v", err)
	}
	assert.Contains(t, result.Stdout, expectedIP, "Expected %s to resolve to %s", host, expectedIP)
}
