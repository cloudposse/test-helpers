package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)


func NewElbV2Client(t *testing.T, region string) *elasticloadbalancingv2.Client {
	client, err := NewElbV2ClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewElbV2ClientE creates an ELB client.
func NewElbV2ClientE(t *testing.T, region string) (*elasticloadbalancingv2.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return elasticloadbalancingv2.NewFromConfig(*sess), nil
}
