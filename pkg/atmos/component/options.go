package component

type Stack struct {
	Component string
	RandomID  string
	Region    string
	StackName string
}

type AwsComponentTestOptions struct {
	AwsRegion         string
	ComponentName     string
	FixturesPath      string
	SkipAwsNuke       bool
	StackDependencies []Stack
	StackName         string
}

// Option type represents a configuration option
type AwsComponentTestOption func(*AwsComponentTestOptions)

// WithFixturesPath is an option for setting the FixturesPath
func WithFixturesPath(fixturesPath string) AwsComponentTestOption {
	return func(a *AwsComponentTestOptions) {
		a.FixturesPath = fixturesPath
	}
}

// WithSkipAwsNuke is an option for setting SkipAwsNuke
func WithSkipAwsNuke(skip bool) AwsComponentTestOption {
	return func(a *AwsComponentTestOptions) {
		a.SkipAwsNuke = skip
	}
}

// WithStackDependencies is an option for setting the StackDependencies
func WithDependencies(dependencies []Stack) AwsComponentTestOption {
	return func(a *AwsComponentTestOptions) {
		a.StackDependencies = dependencies
	}
}

// NewAwsComponentTestOptions creates a new AwsComponentTestOptions with required fields and optional configuration
func NewAwsComponentTestOptions(awsRegion, componentName, stackName string, opts ...AwsComponentTestOption) AwsComponentTestOptions {
	options := &AwsComponentTestOptions{
		AwsRegion:     awsRegion,
		ComponentName: componentName,
		StackName:     stackName,
		FixturesPath:  "./fixtures",
		SkipAwsNuke:   false,
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(options)
	}

	return *options
}
