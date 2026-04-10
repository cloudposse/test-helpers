package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)

func NewMSKClient(t *testing.T, region string) *kafka.Client {
	client, err := NewMSKClientE(t, region)
	require.NoError(t, err)

	return client
}

func NewMSKClientE(t *testing.T, region string) (*kafka.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return kafka.NewFromConfig(*sess), nil
}
