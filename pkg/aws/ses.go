package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)



// NewElbV2Client creates en ELB client.
func NewSESV2Client(t *testing.T, region string) *sesv2.Client {
	client, err := NewSESV2ClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewSESV2ClientE creates an SES v2 client.
func NewSESV2ClientE(t *testing.T, region string) (*sesv2.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return sesv2.NewFromConfig(*sess), nil
}
