package examples_helper

import (
	"os"
	"path/filepath"
	"testing"

	log "github.com/charmbracelet/log"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
	"github.com/stretchr/testify/require"
)

func vendorFileExists(t *testing.T, config *c.Config) bool {
	vendorFile := filepath.Join(config.TempDir, "vendor.yaml")
	_, err := os.Stat(vendorFile)
	if os.IsNotExist(err) {
		log.WithPrefix(t.Name()).Warn("skipping test phase because vendor file does not exist", "phase", "setup test suite/vendor dependencies", "vendorFile", vendorFile)
		return false
	}

	require.NoError(t, err)

	return true
}

func (s *TestSuite) VendorDependencies(t *testing.T, config *c.Config) {
	const phaseName = "vendor dependencies"
	if !vendorFileExists(t, config) {
		return
	}

	if !vendorFileExists(t, config) || config.SkipVendorDependencies {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	log.Debug("running atmos vendor pull in tempdir")
	atmosOptions := getAtmosOptionsFromSetupConfiguration(t, s.Config, s.SetupConfiguration, "", "", nil, nil)
	_, err := atmos.VendorPullE(t, atmosOptions)
	if err != nil {
		s.logPhaseStatus(phaseName, "error")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")
}

func (s *TestSuite) pullVendorYamlComponents(t *testing.T, config *c.Config) {
	const phaseName = "pulling vendor.yaml components"
	if !vendorFileExists(t, config) {
		return
	}

	if !vendorFileExists(t, config) || config.SkipVendorDependencies {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	log.Debug("running atmos vendor pull in tempdir")
	atmosOptions := GetAtmosOptions(t, config, "", "", nil)
	_, err := atmos.VendorPullE(t, atmosOptions)
	if err != nil {
		s.logPhaseStatus(phaseName, "error")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")
}

func (s *TestSuite) pullComponentYamlComponents(t *testing.T, config *c.Config) {
	const phaseName = "pulling component.yaml components"
	s.logPhaseStatus(phaseName, "started")

	filepath.Walk(s.Config.TempDir, func(path string, info os.FileInfo, err error) error {
		log.WithPrefix(t.Name()).WithPrefix("Pulling Component Yaml").Debug("checking", "path", path)
		if err != nil {
			return err
		}

		if filepath.Base(path) == "component.yaml" {
			//Get the folder name of component.yaml
			component := filepath.Base(filepath.Dir(path))

			// Get the atmos options for the component
			atmosOptions := getAtmosOptionsFromSetupConfiguration(t, config, s.SetupConfiguration, component, "", nil, nil)
			log.WithPrefix(s.T().Name()).Info("found component.yaml", "path", path, "component", component)

			// Pull the component
			_, _ = atmos.VendorPullComponent(t, atmosOptions)
		}
		return nil
	})

	s.logPhaseStatus(phaseName, "completed")
}

func (s *TestSuite) VendorComponents(t *testing.T, config *c.Config) {
	s.pullVendorYamlComponents(t, config)
	s.pullComponentYamlComponents(t, config)
}
