package dependency

type Dependency struct {
	AdditionalVars *map[string]interface{}
	ComponentName  string
	StackName      string
}
