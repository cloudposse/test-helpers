package dependency

import "github.com/cloudposse/test-helpers/pkg/atmos"

type Dependency struct {
	AdditionalVars     *map[string]interface{}
	ComponentName      string
	StackName          string
	Function           func() error
	Args               []string
	Vendor             bool
	VendorOnly         bool
	Targets            []string
	AddRandomAttribute bool
	Options            *atmos.Options
	WorkflowName       string
	WorkflowFile       string
}
