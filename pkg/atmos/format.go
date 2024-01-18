package atmos

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/collections"
	tt "github.com/gruntwork-io/terratest/modules/terraform"
)

const terraformCmd = "terraform"

// TerraformCommandsWithPlanFileSupport is a list of all the Terraform commands that support interacting with plan
// files.
var TerraformCommandsWithPlanFileSupport = []string{
	"plan",
	"apply",
	"show",
	"graph",
}

// FormatTerraformArgs converts the inputs to a format palatable to terraform. This includes converting the given vars
// to the format the Terraform CLI expects (-var key=value).
func FormatAtmosTerraformArgs(options *Options, args ...string) []string {
	var terraformArgs []string
	commandType := args[0]

	lockSupported := collections.ListContains(tt.TerraformCommandsWithLockSupport, commandType)
	planFileSupported := collections.ListContains(TerraformCommandsWithPlanFileSupport, commandType)

	terraformArgs = append(terraformArgs, "terraform", commandType, options.Component, "-s", options.Stack)

	// Include -var and -var-file flags unless we're running 'apply' with a plan file
	includeVars := !(commandType == "apply" && len(options.PlanFilePath) > 0)

	terraformArgs = append(terraformArgs, args[1:]...)

	if includeVars {
		if options.SetVarsAfterVarFiles {
			terraformArgs = append(terraformArgs, tt.FormatTerraformArgs("-var-file", options.VarFiles)...)
			terraformArgs = append(terraformArgs, tt.FormatTerraformVarsAsArgs(options.Vars)...)
		} else {
			terraformArgs = append(terraformArgs, tt.FormatTerraformVarsAsArgs(options.Vars)...)
			terraformArgs = append(terraformArgs, tt.FormatTerraformArgs("-var-file", options.VarFiles)...)
		}
	}

	terraformArgs = append(terraformArgs, tt.FormatTerraformArgs("-target", options.Targets)...)

	if options.NoColor {
		terraformArgs = append(terraformArgs, "-no-color")
	}

	if options.RedirectStrErrDestination != "" {
		terraformArgs = append(terraformArgs, fmt.Sprintf("--redirect-stderr=%s", options.RedirectStrErrDestination))
	}

	if lockSupported {
		// If command supports locking, handle lock arguments
		terraformArgs = append(terraformArgs, tt.FormatTerraformLockAsArgs(options.Lock, options.LockTimeout)...)
	}

	if planFileSupported {
		// The plan file arg should be last in the terraformArgs slice. Some commands use it as an input (e.g. show, apply)
		terraformArgs = append(terraformArgs, tt.FormatTerraformPlanFileAsArg(commandType, options.PlanFilePath)...)
	}

	return terraformArgs
}

func FormatArgs(options *Options, args ...string) []string {
	var atmosArgs []string
	commandType := args[0]

	if commandType == terraformCmd {
		atmosArgs = append(atmosArgs, FormatAtmosTerraformArgs(options, args[1:]...)...)
	}

	return atmosArgs
}
