package atmos

type DescribeStacksTerraformComponent struct {
	AtmosComponent         string      `json:"atmos_component"`
	AtmosStack             string      `json:"atmos_stack"`
	AtmosStackFile         string      `json:"atmos_stack_file"`
	Backend                interface{} `json:"backend"`
	BackendType            string      `json:"backend_type"`
	Command                string      `json:"command"`
	Component              string      `json:"component"`
	Env                    interface{} `json:"env"`
	Inheritance            interface{} `json:"inheritance"`
	Metadata               interface{} `json:"metadata"`
	Overrides              interface{} `json:"overrides"`
	RemoteStateBackend     interface{} `json:"remote_state_backend"`
	RemoteStateBackendType string      `json:"remote_state_backend_type"`
	Settings               interface{} `json:"settings"`
	Vars                   interface{} `json:"vars"`
	Workspace              string      `json:"workspace"`
}

type DescribeStacksTerraform struct {
	Terraform map[string]DescribeStacksTerraformComponent `json:"terraform"`
}

type DescribeStacksComponent struct {
	Component map[string]any `json:"component"`
}

type DescribeStacksOutput struct {
	Stacks map[string]DescribeStacksComponent
}
