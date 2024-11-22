

<!-- markdownlint-disable -->
# test-helpers <a href="https://cpco.io/homepage?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/test-helpers&utm_content="><img align="right" src="https://cloudposse.com/logo-300x69.svg" width="150" /></a>
<a href="https://github.com/cloudposse/test-helpers/releases/latest"><img src="https://img.shields.io/github/release/cloudposse/test-helpers.svg" alt="Latest Release"/></a><a href="https://slack.cloudposse.com"><img src="https://slack.cloudposse.com/badge.svg" alt="Slack Community"/></a>
<!-- markdownlint-restore -->

<!--




  ** DO NOT EDIT THIS FILE
  **
  ** This file was automatically generated by the `cloudposse/build-harness`.
  ** 1) Make all changes to `README.yaml`
  ** 2) Run `make init` (you only need to do this once)
  ** 3) Run`make readme` to rebuild this file.
  **
  ** (We maintain HUNDREDS of open source projects. This is how we maintain our sanity.)
  **





-->

`test-helpers` is a library that adds some missing functionality to [terratest](https://terratest.gruntwork.io).




## Introduction


`test-helpers` is a library that adds some missing functionality to [terratest](https://terratest.gruntwork.io).

`test-helpers` includes functionality for:

   - Destroying all resources in an AWS account after a test run using [aws-nuke](https://github.com/ekristen/aws-nuke)
   - Running tests with [atmos](https://github.com/cloudposse/atmos) stack configs
   - Running tests on components (Terraform root modules) that follow the Cloud Posse standards for implementation


## Install

Install the latest version in your go tests

```console
go get github.com/cloudposse/test-helpers
```

Get a specific version

```console
go get github.com/cloudposse/test-helpers@v0.0.1
```

## Usage

You can use `test-helpers` as a library in your own Go test code along with the `terratest` library.

## Packages

### pkg/atmos

This library is designed to be used with [atmos](https://github.com/cloudposse/atmos) to allow you to run tests with
different stack configurations. For example, imagine you have a component that you want to test to make sure it applies
properly and that the coput contains "Hello, World". Below hows how you could run a test from an atmos stack.

```go
func TestApplyNoError(t *testing.T) {
  t.Parallel()

  testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/atmos", t.Name())
  require.NoError(t, err)
  defer os.RemoveAll(testFolder)

  fmt.Printf("running in %s\n", testFolder)

  options := WithDefaultRetryableErrors(t, &Options{
    AtmosBasePath: testFolder,
    Component:     "terraform-no-error",
    Stack:         "test-test-test",
    NoColor:       true,
  })

  out := Apply(t, options)

  require.Contains(t, out, "Hello, World")
}
```

### pkg/aws-nuke

```go
func TestAwsNuke(t *testing.T) {
  t.Parallel()
  randID := strings.ToLower(random.UniqueId())

  rootFolder := "../../"
  terraformFolderRelativeToRoot := "examples/awsnuke-example"
  tempTestFolder := testStructure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)
  defer os.RemoveAll(tempTestFolder)

  terraformOptions := &terraform.Options{
      TerraformDir: tempTestFolder,
      Upgrade:      true,
      VarFiles:     []string{fmt.Sprintf("fixtures.%s.tfvars", testRegion)},
      Vars: map[string]interface{}{
        "attributes": []string{randID},
        "default_tags": map[string]string{
        terratestTagName: randID,
      },
   },
}
```

### pkg/atmos/aws-component-helper

```go
package test

import (
  "fmt"
  "testing"

  "github.com/cloudposse/test-helpers/pkg/atmos"
  helper "github.com/cloudposse/test-helpers/pkg/atmos/aws-component-helper"
  "github.com/stretchr/testify/require"
)

var suite *helper.TestSuite

// TestMain is the entry point for the test suite. It initializes the test
// suite and runs the tests, then tears the test suite down.
func TestMain(m *testing.M) {
  var err error

  // Configure the test suite. The component under test is `bastion` and the
  // stack is `test` in the `us-east-2` region.
  suite, err = helper.NewTestSuite("us-east-2", "bastion", "test")
  if err != nil {
    panic(err)
  }

  // Add dependencies for the component under test in the same stack. If you
  // want to add dependencies in different stacks, use AddCrossStackDependencies.
  //
  // Dependencies are deployed in serial in the order they are added.
  suite.AddDependencies([]string{"vpc"})

  // Create a new testing object since TestMain doesn't have one and we need
  // one to call the Setup and Teardown functions
  t := &testing.T{}

  defer suite.TearDown(t)
  suite.Setup(t)

  m.Run()
}

func TestBastion(t *testing.T) {
  additionalVars := map[string]interface{}{}
  defer suite.DestroyComponentUnderTest(t, additionalVars)

  _, err := suite.DeployComponentUnderTest(t, additionalVars)
  require.NoError(t, err)

  instanceProfile := atmos.Output(t, suite.AtmosOptions, "iam_instance_profile")
  require.Equal(t, instanceProfile, fmt.Sprintf("eg-cptest-ue2-test-bastion-%s", suite.RandomSeed))
}
```

```bash
$ go test -v -run TestBastion -skip-aws-nuke
```

#### Test Phases

The `aws-component-helper` test suite is designed to be used with the `atmos` tool to allow you to test components
that follow the Cloud Posse standards for implementation. The test suite runs several test "phases", which can be
individually skipped as needed. By default, all phases are run. The following phases are run by default:

| Phase | Description |Flag|
| ----- | ----------- |----|
| Setup Test Suite | Bootstraps a temp directory and creates a new test suite or reads in a test suite from `test-suite.json` | N/A |
| Setup Component Under Test | Copies the component from `src/` to the temp dir `components/terraform` | `-skip-setup-cut` |
| Vendor Dependencies | Runs the `atmos vendor pull` command to pull in dependency components | `-skip-vendor` |
| Deploy Dependencies | Runs the `atmos terraform apply` command to deploy the dependency stacks | `-skip-deploy-deps` |
| Verify Enabled Flag | Runs a test to ensure the `enabled` flag results in no resources being created | `-skip-verify-enabled-flag` |
| Deploy Component Under Test | Runs the `atmos terraform apply` command to deploy the component we are testing | `-skip-deploy-cut` |
| Destroy Component Under Test | Runs the `atmos terraform destroy` command to destroy the component we are testing | `-skip-destroy-cut` |
| Destroy Dependencies | Runs the `atmos destroy` command to destroy the dependencies | `-skip-destroy-deps` |
| Tear Down Test Suite | Cleans up the temp directory | `-skip-teardown` |
| Nuke Test Account | Uses [aws-nuke](https://github.com/ekristen/aws-nuke) to destroy all resources created during the test (by tag) | `-skip-aws-nuke` |

## Examples

The [example](examples/) folder contains a full set examples that demonstrate the use of `test-helpers`:

  - [example](examples/awsnuke-example) folder contains a terraform module that can be used to test the `awsnuke` functionality.
  The test for this module is in [pkg/awsnuke/awsnuke_test.go](pkg/awsnuke/awsnuke_test.go).
  - [test/aws-component-helper](test/aws-component-helper) folder contains a terraform module and test that can be used to demonstrate the `aws-component-helper` functionality.













## ✨ Contributing

This project is under active development, and we encourage contributions from our community.



Many thanks to our outstanding contributors:

<a href="https://github.com/cloudposse/test-helpers/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=cloudposse/test-helpers&max=24" />
</a>

For 🐛 bug reports & feature requests, please use the [issue tracker](https://github.com/cloudposse/test-helpers/issues).

In general, PRs are welcome. We follow the typical "fork-and-pull" Git workflow.
 1. Review our [Code of Conduct](https://github.com/cloudposse/test-helpers/?tab=coc-ov-file#code-of-conduct) and [Contributor Guidelines](https://github.com/cloudposse/.github/blob/main/CONTRIBUTING.md).
 2. **Fork** the repo on GitHub
 3. **Clone** the project to your own machine
 4. **Commit** changes to your own branch
 5. **Push** your work back up to your fork
 6. Submit a **Pull Request** so that we can review your changes

**NOTE:** Be sure to merge the latest changes from "upstream" before making a pull request!

### 🌎 Slack Community

Join our [Open Source Community](https://cpco.io/slack?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/test-helpers&utm_content=slack) on Slack. It's **FREE** for everyone! Our "SweetOps" community is where you get to talk with others who share a similar vision for how to rollout and manage infrastructure. This is the best place to talk shop, ask questions, solicit feedback, and work together as a community to build totally *sweet* infrastructure.

### 📰 Newsletter

Sign up for [our newsletter](https://cpco.io/newsletter?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/test-helpers&utm_content=newsletter) and join 3,000+ DevOps engineers, CTOs, and founders who get insider access to the latest DevOps trends, so you can always stay in the know.
Dropped straight into your Inbox every week — and usually a 5-minute read.

### 📆 Office Hours <a href="https://cloudposse.com/office-hours?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/test-helpers&utm_content=office_hours"><img src="https://img.cloudposse.com/fit-in/200x200/https://cloudposse.com/wp-content/uploads/2019/08/Powered-by-Zoom.png" align="right" /></a>

[Join us every Wednesday via Zoom](https://cloudposse.com/office-hours?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/test-helpers&utm_content=office_hours) for your weekly dose of insider DevOps trends, AWS news and Terraform insights, all sourced from our SweetOps community, plus a _live Q&A_ that you can’t find anywhere else.
It's **FREE** for everyone!
## License

<a href="https://opensource.org/licenses/Apache-2.0"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge" alt="License"></a>

<details>
<summary>Preamble to the Apache License, Version 2.0</summary>
<br/>
<br/>

Complete license is available in the [`LICENSE`](LICENSE) file.

```text
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
```
</details>

## Trademarks

All other trademarks referenced herein are the property of their respective owners.


---
Copyright © 2017-2024 [Cloud Posse, LLC](https://cpco.io/copyright)


<a href="https://cloudposse.com/readme/footer/link?utm_source=github&utm_medium=readme&utm_campaign=cloudposse/test-helpers&utm_content=readme_footer_link"><img alt="README footer" src="https://cloudposse.com/readme/footer/img"/></a>

<img alt="Beacon" width="0" src="https://ga-beacon.cloudposse.com/UA-76589703-4/cloudposse/test-helpers?pixel&cs=github&cm=readme&an=test-helpers"/>
