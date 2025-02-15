package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
	"testing"
)

func NewDocDBlient(t *testing.T, region string) *docdb.Client {
	client, err := NewDocDBlientE(t, region)
	require.NoError(t, err)

	return client
}

func NewDocDBlientE(t *testing.T, region string) (*docdb.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return docdb.NewFromConfig(*sess), nil
}
