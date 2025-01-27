# Atmos Component Helper

## Overview

The `component-helper` package (The Helper) is designed to be used to test components that follow the Cloud Posse
convention for Components (terraform root modules) using standard go testing and the `atmos` tool for creating test
fixtures.

The Helper package provides the concept of a `TestSuite`, which is a collection of tests that are run
against a common component and its dependencies.

The Helper provides several standard test phases that are run in order to setup, deploy, and teardown the component
under test and its dependencies.

The main goals of The Helper are to provide a simple way to test components with as little boilerplate code as possible
while also making it easy for non-Go developers to understand and use.

## Example

```go
package test

import (
  "fmt"
  "testing"

  atmos "github.com/cloudposse/test-helpers/pkg/atmos"
  helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
  "github.com/stretchr/testify/assert"
)

// First we define a test suite that will contain all of our tests. Here we embed the TestSuite struct from the Helper
// package so that this test suite struct will inherit all of the Helper's functionality.
type VpcTestSuite struct {
  helper.TestSuite
}

// Next we can define a set of tests that will be run for out component. These tests are defined as methods on the suite
// struct so that we can control the lifecycle of dependencies and the test suite. Any methods that start with `Test`
// are considered tests.
func (s *VpcTestSuite) TestVPC() {

  // The component name and stack name are passed to the DeployAtmosComponent()  and DestroyAtmosComponent() methods.
  // These stacks and components are defined in the `fixtures/` directory (see test phases below for more details).
  const component = "vpc"
  const stack = "test-use2-sandbox"

  // First, we need to call the DestroyAtmosComponent() method to ensure that the component is destroyed after the
  // tests are run, even if the tests fail. Although this seems backwards, it is actually a good practice to ensure
  // that the component is destroyed after the tests are run. And Go will not let you leave a function scope without
  // executing the defer statement, even if other parts of the test suite fail.
  defer s.DestroyAtmosComponent(s.T(), component, stack, nil)

  // Next, we need to call the DeployAtmosComponent() method to actually deploy the component. Under the hood, this
  // method will call the `atmos terraform deploy` command to deploy the component.
  options, _ := s.DeployAtmosComponent(s.T(), component, stack, nil)

  // Finally, as an example, we can use the atmos Output() method to get the value of the `vpc_cidr` output from the
  // component and assert that it is equal to the expected value.
  cidrBlock := atmos.Output(s.T(), options, "vpc_cidr")
  assert.Equal(s.T(), "10.1.0.0/16", cidrBlock)
}

func (s *VpcTestSuite) TestVPC2() {
 // Add other tests as needed. For example, you may deploy one VPC with only private subnets, one without a NAT gateway,
 // etc. and run various other assertions to validate the component behaves as you expect with different configuration.
}

// Cloud Posse strongly recommends that each component have a flag (variable) that can be set to enable ot disable the
// entire component (root module). This is extremely useful in a number of scenarios. As such, all of Cloud Posse's
// components have a `enabled` flag that can be set to `true` or `false`. This test verifies that running
// `atmos terraform plan` on the component results in a `no op` plan (that is, a plan that shows no changes).
func (s *VpcTestSuite) TestEnabledFlag() {
  s.VerifyEnabledFlag("vpc", "test-use2-sandbox", nil)
}


// Finally, we define a standard go test function that will run the test suite. Note that this function is not a method
// on the suite struct, but rather a standard go test function. This becomes the entry point for the test suite.
func TestRunVPCSuite(t *testing.T) {

  // Create a new instance of the test suite struct.
  suite := new(VpcTestSuite)

  // Add dependencies to the test suite.
  suite.AddDependency(t, "vpc-flow-logs", "test-use2-sandbox", nil)

  // Run the test suite
  helper.Run(t, suite)
}
```

Now that we have a TestSuite defined, we can use Go's standard test runner

## Component Directory Structure

The Helper assumes that each test suite is used to test a component (root module) and that the component, atmos stack
configuration, and tests are are located within a directory that is structured as described below, however, The Helper
is configurable so that it can be used in other directory structures.

```text
component-root/
├── src/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── other files...
└── test/
    ├── fixtures/
    │   ├── stacks/
    │   │   ├── orgs/
    │   │   │   └── other stack directories...
    │   │── atmos.yaml
    │   └── vendor.yaml
    │
    ├── component_name_test.go
    ├── component_name_variation_test.go
    └── other test files as needed...
```

## Test Phases

As mentioned above, The Helper provides several standard test phases that are run in order to setup, deploy, and
teardown the component under test and its dependencies. It is possible to skip some of these phases using command-line
flags so that you can iterate locally on the tests without having to run all of the phases. This is particularly helpful
when deploying dependencies, which often can take a considerable amount of time to complete.

The following test phases are run by default (and the flag to skip where applicable):

### Setup (--skip-setup)

During the setup phase, The Helper will first create a temporary directory to use for running the tests. If you want to
specify the path to the temporary directory, you can use the `--temp-dir` flag, otherwise it will be randomly generated.

This temporary directory is structured exactly like a typical atmos project directory:

```text
root/
├── components/
│   ├── terraform/
│   │   ├── target/
│   │   ├── componentX/
│   │   └── componentY/
├── stacks/
├── atmos.yaml
└── vendor.yaml
```

The Helper will copy the component under test from the component's `src` directory to the `components/terraform/target`
directory under the temp directory. The source directory can be overridden using the `--src-dir` flag.

Next, the Helper will copy the contents of the `fixtures` directory to the root of the temp directory. This includes the
`stacks` directory, the `atmos.yaml` file, and the `vendor.yaml` file. If yyou want to override the path to the
`fixtures` directory that is being copied, you can use the `--fixtures-path` flag.

### Vendor Dependencies (--skip-vendor)

During the next phase, The Helper will switch into the temp directory and run `atmos vendor pull` to install any
dependencies that were defined in the `vendor.yaml` file.

### Deploy Dependencies (--skip-deploy)

Next, The Helper will switch into the temp directory and run `atmos deploy` for each of the stack dependencies defined
within the test suite. The order of deployment is the same as the order in which the dependencies are added to the suite.

### Test

The Helper will then use the `go test` command to run any tests that are defined in the test suite.

### Destroy Dependencies (--skip-destroy-dependencies)

Once all of the tests have been run, The Helper will destroy all of the dependencies in the reverse order of which they
were deployed by running the `atmos destroy` command for each dependency.

### Teardown (--skip-teardown)

Finally, The Helper will clean up the temporary directory and any other resources that were created during the test run.

## Advanced Usage

- If you choose to use any of the `--skip-*` flags, the test suite will write a file to the directory where the test suite
  is being run. By default this file is called `test_suite.yaml` and it contains information needed to run the test suite
  again in the future, including the path to the temporary and state directories. If you need to run multiple suites at
  the same time, you can use the `-config` flag to specify a different configuration file name.

- When you use the `--skip-setup` flag, the test suite will use the previously created temporary directory and state
  directories, but will always copy the component under test and the `fixtures` directory to the temporary directory to
  ensure that the test suite is always being run on the latest version of the component and its configuration.

- As previously mentioned, The temporary directory is structured exactly like a typical `atmos` project directory. This
  means that you can switch to the temporary directory and run `atmos` commands directly to deploy and destroy the
  component under test and its dependencies. This can be useful for debugging issues with the component or its
  dependencies.

## Flags reference

| Flag                       | Description                                           | Default           |
| -------------------------- | ----------------------------------------------------- | ----------------- |
| -config                    | The path to the config file                           | test_suite.yaml   |
| -fixtures-dir              | The path to the fixtures directory                    | fixtures          |
| -skip-deploy-dependencies  | Skips running the deploy dependencies phase of tests  | false             |
| -skip-destroy-component    | Skips running the destroy component phase of tests    | false             |
| -skip-destroy-dependencies | Skips running the destroy dependencies phase of tests | false             |
| -skip-enabled-flag-test    | Skips running the Enabled flag test                   | false             |
| -skip-setup                | Skips running the setup test suite phase of tests     | false             |
| -skip-teardown             | Skips running the teardown test suite phase of tests  | false             |
| -skip-vendor               | Skips running the vendor dependencies phase of tests  | false             |
| -src-dir                   | The path to the component source directory            | src               |
| -state-dir                 | The path to the terraform state directory             | {temp_dir}/state  |
| -temp-dir                  | The path to the temp directory                        | {random temp dir} |
