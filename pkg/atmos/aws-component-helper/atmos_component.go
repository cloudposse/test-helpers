package aws_component_helper

import (
	"strings"

	"github.com/gruntwork-io/terratest/modules/random"
)

type Output struct {
	Sensitive bool        `json:"sensitive"`
	Value     interface{} `json:"value"`
}

type AtmosComponent struct {
	RandomIdentifier string
	ComponentName    string
	StackName        string
	Vars             map[string]interface{}
	output           map[string]Output
}

func NewAtmosComponent(component string, stack string, vars map[string]interface{}) *AtmosComponent {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)
	return &AtmosComponent{
		RandomIdentifier: randomId,
		ComponentName:    component,
		StackName:        stack,
		Vars:             vars,
	}
}

// func (ac *AtmosComponent) getAtmosOptions(options *atmos.Options, vars map[string]interface{}) *atmos.Options {
// 	result := &atmos.Options{}
// 	if options != nil {
// 		result, _ = options.Clone()
// 	}

// 	// currentTFDataDir := ".terraform"
// 	// if value, ok := options.EnvVars["TF_DATA_DIR"]; ok {
// 	// 	currentTFDataDir = value
// 	// }
// 	// stack := strings.Replace(ac.StackName, "/", "-", -1)
// 	// name := strings.Replace(ac.ComponentName, "/", "-", -1)
// 	// envvars := map[string]string{
// 	// 	// We need to split the TF_DATA_DIR for parallel suites mode
// 	// 	"TF_DATA_DIR":             filepath.Join(currentTFDataDir, fmt.Sprintf("component-%s", ac.RandomIdentifier)),
// 	// 	"TEST_WORKSPACE_TEMPLATE": fmt.Sprintf("%s-%s-%s", stack, name, ac.RandomIdentifier),
// 	// }

// 	// err := mergo.Merge(&result.Vars, ac.vars)
// 	// require.NoError(t, err)

// 	// // Merge in any additional vars passed in
// 	// err = mergo.Merge(&result.Vars, vars)
// 	// require.NoError(t, err)

// 	// result.Component = ac.ComponentName
// 	// result.Stack = ac.StackName

// 	atmosOptions := atmos.WithDefaultRetryableErrors(, result)
// 	return atmosOptions
// }
