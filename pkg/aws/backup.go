package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/stretchr/testify/require"
	"github.com/gruntwork-io/terratest/modules/aws"
)

func NewBackupClient(t *testing.T, region string) *backup.Client {
	client, err := NewBackupClientE(t, region)
	require.NoError(t, err)

	return client
}

func NewBackupClientE(t *testing.T, region string) (*backup.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return backup.NewFromConfig(*sess), nil
}
