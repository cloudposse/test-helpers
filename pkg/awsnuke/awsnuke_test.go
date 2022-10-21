package awsnuke

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	testStructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

const testRegion string = "us-east-2"
const terratestTagName string = "CreatedByTerratestRun"

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

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Assert that the bucket has been created
	bucket, err := aws.FindS3BucketWithTagE(t, testRegion, terratestTagName, randID)
	assert.NoError(t, err)
	assert.Equal(t, bucket, fmt.Sprintf("eg-test-s3-bucket-test-%s", randID))

	// Nuke the account with our config
	NukeTestAccountByTag(t, terratestTagName, randID, []string{testRegion}, false)

	// Assert that the bucket doesn't exist anymore
	bucket, err = aws.FindS3BucketWithTagE(t, testRegion, terratestTagName, randID)
	assert.NoError(t, err)
	assert.Empty(t, bucket)
}
