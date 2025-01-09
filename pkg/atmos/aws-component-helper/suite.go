package aws_component_helper

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Flags for controlling test behavior
var (
	skipTests = flag.Bool("skip-tests", false, "skip tests")
	skipSetup = flag.Bool("skip-setup", false, "skip setup")
)

// Structure to define a teardown step
type tearDown struct {
	component *AtmosComponent                 // The component to be torn down
	callback  *func(t *testing.T, atm *Atmos) // Optional teardown callback
}

// Suite represents a test suite and its associated configurations, state, and dependencies
type Suite struct {
	t                *testing.T        // Testing object
	randomIdentifier string            // Unique identifier for the suite
	name             string            // Suite name
	dependencies     []*AtmosComponent // Dependencies for the suite
	teardown         []*tearDown       // Teardown steps
	options          *atmos.Options    // Atmos options for the suite
	state            *State            // State associated with the suite
}

// NewSuite initializes a new Suite instance
func NewSuite(t *testing.T, name string, fixture *Fixture) *Suite {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	// Ensure suite name does not contain invalid characters
	require.NotContains(t, name, "/")

	// Fork the state for the suite
	suiteState, err := fixture.State.Fork(name)
	require.NoError(t, err)

	// Initialize and return the Suite
	return &Suite{
		t:                t,
		name:             name,
		randomIdentifier: randomId,
		dependencies:     []*AtmosComponent{},
		state:            suiteState,
		teardown:         []*tearDown{},
		options:          fixture.getAtmosOptions(&atmos.Options{}, map[string]interface{}{}),
	}
}

// AddDependency adds a component dependency to the suite and deploys it
func (ts *Suite) AddDependency(componentName string, stackName string) {
	component := NewAtmosComponent(componentName, stackName, nil)
	ts.dependencies = append(ts.dependencies, component)
	ts.teardown = append(ts.teardown, &tearDown{component: component, callback: nil})

	if *skipSetup {
		fmt.Printf("Skip suite %s setup dependency component: %s stack: %s\n", ts.name, componentName, stackName)
		return
	}

	// Deploy the component using Atmos with suite TF state
	ts.getAtmos(ts.state).Deploy(component)
}

// getAtmos returns an Atmos instance configured for the suite's state
func (ts *Suite) getAtmos(state *State) *Atmos {
	return NewAtmos(ts.t, state, ts.getAtmosOptions(map[string]interface{}{
		"attributes": []string{ts.GetRandomIdentifier()},
	}))
}

// getTestAtmos returns an Atmos instance for test-specific purposes
func (ts *Suite) getTestAtmos(state *State) *Atmos {
	return NewAtmos(ts.t, state, ts.getAtmosOptions(map[string]interface{}{}))
}

// runTeardown executes all teardown steps for the suite
func (ts *Suite) runTeardown() {
	if *skipTeardown {
		fmt.Printf("Skip teardown suite %s\n", ts.name)
		return
	}

	atm := ts.getAtmos(ts.state)
	for i := len(ts.teardown) - 1; i >= 0; i-- {
		step := ts.teardown[i]

		// Execute the callback, if provided
		if step.callback != nil {
			(*step.callback)(ts.t, atm)
		}

		// Destroy the component, if applicable
		if step.component != nil {
			atm.Destroy(step.component)
		}
	}

	// Teardown the suite's state
	err := ts.state.Teardown()
	assert.NoError(ts.t, err)
}

// Setup executes a setup callback for the suite
func (ts *Suite) Setup(t *testing.T, f func(t *testing.T, atm *Atmos)) {
	if *skipSetup {
		fmt.Printf("Skip suite %s setup callback\n", ts.name)
		return
	}
	// Get Atmos with suite TF state
	atm := ts.getAtmos(ts.state)
	f(t, atm)
}

// TearDown adds a custom teardown callback to the suite
func (ts *Suite) TearDown(t *testing.T, f func(t *testing.T, atm *Atmos)) {
	ts.teardown = append(ts.teardown, &tearDown{component: nil, callback: &f})
}

// Test runs a test within the suite
func (ts *Suite) Test(t *testing.T, name string, f func(t *testing.T, atm *Atmos)) {
	if *skipTests {
		fmt.Printf("Skip test %s/%s\n", ts.name, name)
		return
	}

	// Fork the state for the test
	testState, err := ts.state.Fork(name)
	require.NoError(t, err)
	defer testState.Teardown() // Ensure test state is torn down after execution

	// Get Atmos with test TF state
	atm := ts.getTestAtmos(testState)

	// Run the test if it matches the filter
	testRunName := fmt.Sprintf("%s/%s", t.Name(), name)
	if ok, err := matchFilter(testRunName); ok {
		t.Run(name, func(t *testing.T) {
			f(t, atm)
		})
	} else {
		require.NoError(t, err)
	}
}

// getAtmosOptions generates Atmos options for the suite with additional variables
func (ts *Suite) getAtmosOptions(vars map[string]interface{}) *atmos.Options {
	result, err := ts.options.Clone()
	require.NoError(ts.t, err)

	// Merge the provided variables into the options
	err = mergo.Merge(&result.Vars, vars)
	require.NoError(ts.t, err)

	return result
}

// GetRandomIdentifier returns the suite's unique random identifier
func (ts *Suite) GetRandomIdentifier() string {
	if *devMode {
		return ts.state.GetIdentifier()
	}
	return ts.randomIdentifier
}
