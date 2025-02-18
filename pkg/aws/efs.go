package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)


func NewEFSClient(t *testing.T, region string) *efs.Client {
	client, err := NewEFSClientE(t, region)
	require.NoError(t, err)

	return client
}

func NewEFSClientE(t *testing.T, region string) (*efs.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return efs.NewFromConfig(*sess), nil
}
