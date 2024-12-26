package aws_component_helper

import (
	"dario.cat/mergo"
	"flag"
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
	"testing"
)

var (
	skipDeploySuiteDependencies  = flag.Bool("cth.skip-deploy-suite-deps", false, "skip deploy suite deps")
	skipDestroySuiteDependencies = flag.Bool("cth.skip-destroy-suite-deps", false, "skip destroy suite deps")
)

type XTestSuite struct {
	RandomIdentifier string
	setup            []*AtmosComponent
	tests            map[string]*ComponentTest
	atmosOptions     *atmos.Options
}

func NewXTestSuite() *XTestSuite {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	return &XTestSuite{
		RandomIdentifier: randomId,
		setup:            make([]*AtmosComponent, 0),
		tests:            make(map[string]*ComponentTest),
	}
}

func (ts *XTestSuite) getAtmosOptions(t *testing.T, options *atmos.Options, vars map[string]interface{}) *atmos.Options {
	result := &atmos.Options{}
	if options != nil {
		result, _ = options.Clone()
	}

	currentTFDataDir := ".terraform"
	if value, ok := options.EnvVars["TF_DATA_DIR"]; ok {
		currentTFDataDir = value
	}

	envvars := map[string]string{
		// We need to split the TF_DATA_DIR for parallel suites mode
		"TF_DATA_DIR": filepath.Join(currentTFDataDir, fmt.Sprintf("suite-%s", ts.RandomIdentifier)),
	}

	err := mergo.Merge(&result.EnvVars, envvars)
	require.NoError(t, err)

	err = mergo.Merge(&result.Vars, vars)
	require.NoError(t, err)

	suiteVars := map[string]interface{}{
		"attributes": []string{ts.RandomIdentifier},
		"default_tags": map[string]string{
			"CreatedByAtmosTestSuite": ts.RandomIdentifier,
		},
	}

	err = mergo.Merge(&result.Vars, suiteVars)
	require.NoError(t, err)
	return result
}
