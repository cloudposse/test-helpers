package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)


func NewNetworkmanagerClient(t *testing.T, region string) *networkmanager.Client {
	client, err := NewNetworkmanagerClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewNetworkmanagerClientE creates a Network Manager client.
func NewNetworkmanagerClientE(t *testing.T, region string) (*networkmanager.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return networkmanager.NewFromConfig(*sess), nil
}
