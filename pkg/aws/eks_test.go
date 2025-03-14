package aws

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"	
	"github.com/gruntwork-io/terratest/modules/retry"
	terratestAWS "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
)

func TestEksCluster(t *testing.T) {
	t.Skip("Skipping EKS cluster test")
	t.Parallel()

	ctx := context.Background()

	// Setting the cluster to always run in us-east-1 and only the us-east-1a, us-east-1b, and us-east-1d subnets as there
	// are seemingly randon subnet restrictions for EKS clusters.
	region := terratestAWS.GetRandomStableRegion(t, []string{"us-east-1"}, nil)
	vpc, err := terratestAWS.GetDefaultVpcE(t, region)
	assert.Nil(t, err)

	var subnetList []string
	for _, subnet := range vpc.Subnets {
		if subnet.AvailabilityZone == "us-east-1a" || subnet.AvailabilityZone == "us-east-1b" || subnet.AvailabilityZone == "us-east-1d" {
			subnetList = append(subnetList, subnet.Id)
		}
	}

	iamClient, err := terratestAWS.NewIamClientE(t, region)
	assert.Nil(t, err)

	role, err := iamClient.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String("AWSServiceRoleForOrganizations"),
	})
	assert.Nil(t, err)

	clusterName := "terratest"
	c1, err := CreateEksClusterE(t, ctx, region, clusterName, *role.Role.Arn, true, true, subnetList, []string{"0.0.0.0/0"}, []string{})
	defer DeleteEksCluster(t, ctx, region, c1)

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *c1.Name)

	maxRetries := 60
	sleepBetweenRetries := 30 * time.Second

	status := retry.DoWithRetry(t, "Ensure cluster is active", maxRetries, sleepBetweenRetries, func() (string, error) {
		cluster := GetEksCluster(t, ctx, region, clusterName)
		status := cluster.Status

		if status != types.ClusterStatusActive {
			return "", fmt.Errorf("Got Cluster Status %s. Retrying.\n", cluster.Status)
		}
		return string(status), nil
	})

	assert.Equal(t, "ACTIVE", status)

	c2, err := GetEksClusterE(t, ctx, region, *c1.Name)

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *c2.Name)
}