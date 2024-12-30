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
	randomIdentifier string
	ComponentName    string
	StackName        string
	Vars             map[string]interface{}
	output           map[string]Output
}

func NewAtmosComponent(component string, stack string, vars map[string]interface{}) *AtmosComponent {
	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	return &AtmosComponent{
		randomIdentifier: randomId,
		ComponentName:    component,
		StackName:        stack,
		Vars:             vars,
	}
}

func (ts *AtmosComponent) GetRandomIdentifier() string {
	return ts.randomIdentifier
}
