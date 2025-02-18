package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)

func NewCloudTrailClient(t *testing.T, region string) *cloudtrail.Client {
	client, err := NewCloudTrailClientE(t, region)
	require.NoError(t, err)

	return client
}

func NewCloudTrailClientE(t *testing.T, region string) (*cloudtrail.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return cloudtrail.NewFromConfig(*sess), nil
}
