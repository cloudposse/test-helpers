package dependency

type Dependency struct {
	AdditionalVars *map[string]interface{}
	ComponentName  string
	StackName      string
	Function       func() error
	Args           []string
}
