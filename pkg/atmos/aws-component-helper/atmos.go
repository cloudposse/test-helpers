package aws_component_helper

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	atmosApply         = atmos.Apply
	atmosDestroy       = atmos.Destroy
	atmosPlanExitCodeE = atmos.PlanExitCodeE
	atmosVendorPull    = atmos.VendorPull
	atmosOutputAllE    = atmos.OutputStructE
)
var (
	skipDeploy  = flag.Bool("skip-deploy", false, "skip all deployments")
	skipDestroy = flag.Bool("skip-destroy", false, "skip all destroy")
)

type Atmos struct {
	t       *testing.T
	options *atmos.Options
	state   *State
}

func NewAtmos(t *testing.T, state *State, options *atmos.Options) *Atmos {
	return &Atmos{
		t:       t,
		options: options,
		state:   state,
	}
}

func (ts *Atmos) GetAndDeploy(componentName string, stackName string, vars map[string]interface{}) *AtmosComponent {
	component := NewAtmosComponent(componentName, stackName, vars)
	ts.Deploy(component)
	return component
}

func (ts *Atmos) GetAndDestroy(componentName string, stackName string, vars map[string]interface{}) *AtmosComponent {
	component := NewAtmosComponent(componentName, stackName, vars)
	ts.Destroy(component)
	return component
}

func (ts *Atmos) Deploy(component *AtmosComponent) {
	options := ts.getAtmosOptions(component)
	defer os.RemoveAll(options.AtmosBasePath)
	err := copyDirectoryRecursively(ts.options.AtmosBasePath, options.AtmosBasePath)
	require.NoError(ts.t, err)
	if !*skipDeploy {
		atmosApply(ts.t, options)
		err := atmosOutputAllE(ts.t, options, "", &component.output)
		require.NoError(ts.t, err)
	}
}

func (ts *Atmos) Destroy(component *AtmosComponent) {
	options := ts.getAtmosOptions(component)
	defer os.RemoveAll(options.AtmosBasePath)
	err := copyDirectoryRecursively(ts.options.AtmosBasePath, options.AtmosBasePath)
	assert.NoError(ts.t, err)
	if !*skipDestroy {
		atmosDestroy(ts.t, options)
	}
}

func (ts *Atmos) loadOutputAll(component *AtmosComponent) {
	if component.output != nil {
		return
	}
	options := ts.getAtmosOptions(component)
	defer os.RemoveAll(options.AtmosBasePath)
	err := copyDirectoryRecursively(ts.options.AtmosBasePath, options.AtmosBasePath)
	require.NoError(ts.t, err)

	err = atmosOutputAllE(ts.t, options, "", &component.output)
	if err != nil && strings.Contains(err.Error(), "Backend initialization required") {
		// Run 'terraform workspace' instead of 'terraform init' as it also select the workspace
		// So terraform output will not fail with "Switch to workspace" json parse error
		_, err := atmos.RunAtmosCommandE(ts.t, options, atmos.FormatArgs(options, "terraform", "workspace")...)
		require.NoError(ts.t, err)
		err = atmosOutputAllE(ts.t, options, "", &component.output)
		require.NoError(ts.t, err)
	}
}

func (ts *Atmos) OutputAll(component *AtmosComponent) map[string]Output {
	ts.loadOutputAll(component)
	return component.output
}

func (ts *Atmos) Output(component *AtmosComponent, key string) string {
	ts.loadOutputAll(component)

	if value, ok := component.output[key]; ok {
		return value.Value.(string)
	}
	require.Fail(ts.t, fmt.Sprintf("Output key %s not found", key))
	return ""
}

func (ts *Atmos) OutputList(component *AtmosComponent, key string) []string {
	ts.loadOutputAll(component)
	if value, ok := component.output[key]; ok {
		if outputList, isList := value.Value.([]interface{}); isList {
			result, err := parseListOutputTerraform(outputList)
			require.NoError(ts.t, err)
			return result
		}
		error := atmos.UnexpectedOutputType{Key: key, ExpectedType: "map or list", ActualType: reflect.TypeOf(value).String()}
		require.Fail(ts.t, error.Error())

	} else {
		require.Fail(ts.t, fmt.Sprintf("Output key %s not found", key))
	}
	return []string{}
}

func (ts *Atmos) OutputMapOfObjects(component *AtmosComponent, key string) map[string]interface{} {
	ts.loadOutputAll(component)
	if value, ok := component.output[key]; ok {
		if outputMap, isMap := value.Value.(map[string]interface{}); isMap {
			return outputMap
		}
		error := atmos.UnexpectedOutputType{Key: key, ExpectedType: "map of objects", ActualType: reflect.TypeOf(value).String()}
		require.Fail(ts.t, error.Error())

	} else {
		require.Fail(ts.t, fmt.Sprintf("Output key %s not found", key))
	}
	return map[string]interface{}{}
}

func (ts *Atmos) OutputStruct(component *AtmosComponent, key string, v any) {
	ts.loadOutputAll(component)
	if value, ok := component.output[key]; ok {
		jsonByte, err := json.Marshal(value.Value)
		require.NoError(ts.t, err)
		jsonString := cleanOutput(string(jsonByte))
		err = json.Unmarshal([]byte(jsonString), &v)
		require.NoError(ts.t, err)
	} else {
		require.Fail(ts.t, fmt.Sprintf("Output key %s not found", key))
	}
}

func (ts *Atmos) getAtmosOptions(component *AtmosComponent) *atmos.Options {
	result, err := ts.options.Clone()
	require.NoError(ts.t, err)

	result.Component = component.ComponentName
	result.Stack = component.StackName

	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	basePath := filepath.Dir(filepath.Clean(ts.options.AtmosBasePath))
	dirName := filepath.Base(ts.options.AtmosBasePath)
	tmpDir := filepath.Join(basePath, fmt.Sprintf(".%s-%s", dirName, randomId))

	result.AtmosBasePath = tmpDir
	resultEnvVars := result.EnvVars
	envvars := map[string]string{
		"ATMOS_BASE_PATH":       result.AtmosBasePath,
		"ATMOS_CLI_CONFIG_PATH": result.AtmosBasePath,
		"TEST_SUITE_NAME":       ts.state.NamespaceDir(),
		"TEST_STATE_DIR":        ts.state.BaseDir(),
	}

	err = mergo.Merge(&envvars, resultEnvVars)
	require.NoError(ts.t, err)

	result.EnvVars = envvars

	if _, ok := result.Vars["attributes"]; !ok {
		result.Vars["attributes"] = []string{component.randomIdentifier}
	}

	if component.Vars != nil {
		err = mergo.Merge(&result.Vars, component.Vars)
		require.NoError(ts.t, err)
	}

	atmosOptions := atmos.WithDefaultRetryableErrors(ts.t, result)

	return atmosOptions
}

func parseMap(m map[string]interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	for k, v := range m {
		switch vt := v.(type) {
		case map[string]interface{}:
			nestedMap, err := parseMap(vt)
			if err != nil {
				return nil, err
			}
			result[k] = nestedMap
		case []interface{}:
			nestedList, err := parseListOfMaps(vt)
			if err != nil {
				return nil, err
			}
			result[k] = nestedList
		case float64:
			testInt, err := strconv.ParseInt((fmt.Sprintf("%v", vt)), 10, 0)
			if err == nil {
				result[k] = int(testInt)
			} else {
				result[k] = vt
			}
		default:
			result[k] = vt
		}

	}
	return result, nil
}

// parseListOfMaps takes a list of maps and parses the types.
// It is mainly a wrapper for parseMap to support lists.
func parseListOfMaps(l []interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	for _, v := range l {

		asMap, isMap := v.(map[string]interface{})
		if !isMap {
			err := errors.New("Type switching to map[string]interface{} failed.")
			return nil, err
		}

		m, err := parseMap(asMap)

		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil

}

func cleanOutput(out string) string {
	var result []rune
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "INFO") {
			continue
		}
		for _, r := range line {
			if r >= 32 && r < 127 { // Keep printable ASCII characters only
				result = append(result, r)
			}
		}
	}
	return string(result)
}
