package atmos

// https://www.terraform.io/docs/commands/plan.html#detailed-exitcode

// TerraformPlanChangesPresentExitCode is the exit code returned by atmos terraform plan detailed exitcode when changes
// are present
const TerraformPlanChangesPresentExitCode = 2

// DefaultSuccessExitCode is the exit code returned when atmos command succeeds
const DefaultSuccessExitCode = 0

// DefaultErrorExitCode is the exit code returned when atmos command fails
const DefaultErrorExitCode = 1

const atmosExamplePath = "../../test/fixtures/atmos"
const testStack = "test-test-test"
