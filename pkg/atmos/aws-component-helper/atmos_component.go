package aws_component_helper

import (
	"strings"

	"github.com/gruntwork-io/terratest/modules/random"
)

// Output represents the structure of a Terraform output with optional sensitivity
type Output struct {
	Sensitive bool        `json:"sensitive"` // Indicates if the output is sensitive
	Value     interface{} `json:"value"`     // The actual value of the output
}

// AtmosComponent represents a Terraform component managed by Atmos
type AtmosComponent struct {
	randomIdentifier string                 // A unique identifier for the component
	ComponentName    string                 // Name of the Terraform component
	StackName        string                 // Name of the stack in which the component resides
	Vars             map[string]interface{} // Variables specific to the component
	output           map[string]Output      // Holds Terraform outputs for the component
}

// Constructor for AtmosComponent
func NewAtmosComponent(component string, stack string, vars map[string]interface{}) *AtmosComponent {
	// Generate a unique random identifier for the component
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	return &AtmosComponent{
		randomIdentifier: randomId,
		ComponentName:    component,
		StackName:        stack,
		Vars:             vars,
	}
}

// GetRandomIdentifier retrieves the unique identifier of the component
func (ts *AtmosComponent) GetRandomIdentifier() string {
	return ts.randomIdentifier
}
