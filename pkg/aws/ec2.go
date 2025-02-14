package aws

import (
	"context"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gruntwork-io/terratest/modules/aws"	
	"github.com/stretchr/testify/assert"
)

func GetNatGatewaysByVpcIdE(t *testing.T, ctx context.Context, vpcId string, awsRegion string) ([]types.NatGateway, error) {
	client, err := aws.NewEc2ClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	filter := &types.Filter{Name: awstypes.String("vpc-id"), Values: []string{vpcId}}
	response, err := client.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{
		Filter: []types.Filter{*filter},
	})
	if err != nil {
		return nil, err
	}
	return response.NatGateways, nil
}

// GetPrivateIpsOfEc2InstancesE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetEc2Instances(t *testing.T, ctx context.Context, instanceID string, awsRegion string) types.Instance {
	ec2Client := aws.NewEc2Client(t, awsRegion)
	// TODO: implement pagination for cases that extend beyond limit (1000 instances)
	input := ec2.DescribeInstancesInput{InstanceIds: []string{instanceID}}
	output, err := ec2Client.DescribeInstances(ctx, &input)
	assert.NoError(t, err)

	return output.Reservations[0].Instances[0]
}
