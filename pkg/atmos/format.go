package atmos

import (
	"fmt"
	"strings"

	"github.com/gruntwork-io/terratest/modules/collections"
	tt "github.com/gruntwork-io/terratest/modules/terraform"
)

const terraformCmd = "terraform"
const vendorCmd = "vendor"
const workflowCmd = "workflow"

// TerraformCommandsWithPlanFileSupport is a list of all the Terraform commands that support interacting with plan
// files.
var TerraformCommandsWithPlanFileSupport = []string{
	"plan",
	"apply",
	"show",
	"graph",
}

func FormatTerraformPlanFileAsArg(commandType string, outPath string) []string {
	if outPath == "" {
		return nil
	}
	if commandType == "plan" {
		return []string{fmt.Sprintf("%s=%s", "-out", outPath)}
	}

	if commandType == "apply" {
		return []string{fmt.Sprintf("%s %s", "--planfile", outPath)}
	}
	return []string{outPath}
}

// FormatTerraformArgs converts the inputs to a format palatable to terraform. This includes converting the given vars
// to the format the Terraform CLI expects (-var key=value).
func FormatAtmosTerraformArgs(options *Options, args ...string) []string {
	var terraformArgs []string
	commandType := args[0]

	lockSupported := collections.ListContains(tt.TerraformCommandsWithLockSupport, commandType)
	planFileSupported := collections.ListContains(TerraformCommandsWithPlanFileSupport, commandType)
	planFileSpecified := len(options.PlanFilePath) > 0
	//skipInit := map[bool]string{true: "--skip-init", false: ""}[options.SkipInit]

	terraformArgs = append(terraformArgs, "terraform", commandType, options.Component, "-s", options.Stack)

	// Include -var and -var-file flags unless we're running 'apply' with a plan file
	includeVars := !(commandType == "apply" && planFileSpecified)

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

	if commandType == "init" {
		terraformArgs = append(terraformArgs, tt.FormatTerraformBackendConfigAsArgs(options.BackendConfig)...)
	}

	if options.NoColor {
		terraformArgs = append(terraformArgs, "-no-color")
	}

	if options.RedirectStrErrDestination != "" {
		terraformArgs = append(terraformArgs, fmt.Sprintf("--redirect-stderr=%s", options.RedirectStrErrDestination))
	}

	if !options.GenerateBackend {
		terraformArgs = append(terraformArgs, fmt.Sprintf("--auto-generate-backend-file=%t", options.GenerateBackend))
	}

	if options.Reconfigure {
		terraformArgs = append(terraformArgs, "-reconfigure")
	}

	if !options.InitRunReconfigure {
		terraformArgs = append(terraformArgs, "--init-run-reconfigure=false")
	}

	if lockSupported {
		// If command supports locking, handle lock arguments
		terraformArgs = append(terraformArgs, tt.FormatTerraformLockAsArgs(options.Lock, options.LockTimeout)...)
	}

	if planFileSupported {
		// The plan file arg should be last in the terraformArgs slice. Some commands use it as an input (e.g. show, apply)
		terraformArgs = append(terraformArgs, FormatTerraformPlanFileAsArg(commandType, options.PlanFilePath)...)
	}

	return terraformArgs
}

// FormatAtmosVendorArgs converts the inputs to a format palatable to atmos vendor.
func FormatAtmosVendorArgs(options *Options, args ...string) []string {
	var vendorArgs []string
	commandType := args[0]

	vendorArgs = append(vendorArgs, "vendor", commandType)

	vendorArgs = append(vendorArgs, args[1:]...)

	if options.RedirectStrErrDestination != "" {
		vendorArgs = append(vendorArgs, fmt.Sprintf("--redirect-stderr=%s", options.RedirectStrErrDestination))
	}

	if options.VendorComponent != "" {
		vendorArgs = append(vendorArgs, "--component", options.VendorComponent)
	}

	if options.VendorStack != "" {
		vendorArgs = append(vendorArgs, "--stack", options.VendorStack)
	}

	if len(options.VendorTags) > 0 {
		vendorArgs = append(vendorArgs, "--tags", strings.Join(options.VendorTags, ","))
	}

	if options.VendorType != "" {
		vendorArgs = append(vendorArgs, "--type", options.VendorType)
	}

	return vendorArgs
}

func FormatArgs(options *Options, args ...string) []string {
	var atmosArgs []string
	commandType := args[0]

	if commandType == terraformCmd {
		atmosArgs = append(atmosArgs, FormatAtmosTerraformArgs(options, args[1:]...)...)
	}

	if commandType == vendorCmd {
		atmosArgs = append(atmosArgs, FormatAtmosVendorArgs(options, args[1:]...)...)
	}

	return atmosArgs
}
