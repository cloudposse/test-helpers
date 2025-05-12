package examples_helper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	log "github.com/charmbracelet/log"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) copyDirectoryContents(srcDir string, destDir string) error {

	// Walk through all files and directories in srcDir and copy them to destDir
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path from srcDir
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Calculate destination path in destDir
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory in destination
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// Copy file content
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			return os.WriteFile(destPath, content, info.Mode())
		}
	})
}

func (s *TestSuite) copyDirectoryRecursively(t *testing.T, srcDir string, destDir string) {
	_, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		path := fmt.Sprintf("setup/copy %s to %s, source directory does not exist", srcDir, destDir)
		s.logPhaseStatus(path, "skipped")
		return
	}

	_, err = os.Stat(destDir)
	if os.IsNotExist(err) {
		path := fmt.Sprintf("setup/copy %s to %s, destination directory does not exist", srcDir, destDir)
		s.logPhaseStatus(path, "skipped")
		return
	}

	log.Debug("copying directory recursively", "srcDir", srcDir, "destDir", destDir)

	err = s.copyDirectoryContents(srcDir, destDir)
	require.NoError(t, err)
}

func setTempDir(t *testing.T, config *c.Config) {
	if config.TempDir == "" {
		tempDir, err := os.MkdirTemp("", "atmos-test-helper")
		require.NoError(t, err)

		viper.Set("TempDir", tempDir)
		config.TempDir = tempDir

		err = config.WriteConfig()
		require.NoError(t, err)
	}

	log.WithPrefix(t.Name()).Info("tests will be run in temp directory", "path", config.TempDir)
}

func setStateDir(t *testing.T, config *c.Config) {
	if config.StateDir == "" {
		stateDir := filepath.Join(config.TempDir, "state")
		viper.Set("StateDir", stateDir)
		config.StateDir = stateDir

		err := os.MkdirAll(stateDir, 0755)
		require.NoError(t, err)

		err = config.WriteConfig()
		require.NoError(t, err)
	}

	log.WithPrefix(t.Name()).Info("terraform state for tests will be saved in state directory", "path", config.StateDir)
}

func setAtmosPaths(t *testing.T, config *c.Config) {
	componentsDir := filepath.Join(config.TempDir, "components", "terraform")
	log.WithPrefix(t.Name()).Debug("creating atmos terraform components directory", "path", componentsDir)
	err := os.MkdirAll(componentsDir, 0755)
	require.NoError(t, err)

	stacksDir := filepath.Join(config.TempDir, "stacks")
	log.WithPrefix(t.Name()).Debug("creating atmos terraform stacks directory", "path", stacksDir)
	err = os.MkdirAll(stacksDir, 0755)
	require.NoError(t, err)
}

func (s *TestSuite) BootstrapTempDir(t *testing.T, config *c.Config) {
	if s.Config.SkipSetupTestSuite {
		return
	}

	setTempDir(t, config)
	setStateDir(t, config)
	setAtmosPaths(t, config)

	s.logPhaseStatus("setup/bootstrap temp dir", "completed")
}

func (s *TestSuite) CopyComponentToTempDir(t *testing.T, config *c.Config) {
	s.logPhaseStatus("setup/copy component to temp dir", "started")

	destPath := filepath.Join(config.TempDir, "components", "terraform", "target")
	err := s.copyDirectoryContents(config.SrcDir, destPath)
	if err != nil {
		s.logPhaseStatus("setup/copy component to temp dir", "failed")
		require.NoError(t, err)
	}

	s.logPhaseStatus("setup/copy component to temp dir", "completed")

	if s.Config.SkipSetupTestSuite {
		return
	}
}

func (s *TestSuite) CopyExampleToTempDir(t *testing.T, config *c.Config) {
	s.logPhaseStatus("setup/copy component to temp dir", "started")

	destPath := filepath.Join(config.TempDir, "components", "terraform", "target")
	err := s.copyDirectoryContents(config.SrcDir, destPath)
	if err != nil {
		s.logPhaseStatus("setup/copy component to temp dir", "failed")
		require.NoError(t, err)
	}

	s.logPhaseStatus("setup/copy component to temp dir", "completed")

	if s.Config.SkipSetupTestSuite {
		return
	}
}
