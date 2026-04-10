package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)

func NewElasticacheClient(t *testing.T, region string) *elasticache.Client {
	client, err := NewElasticacheClientE(t, region)
	require.NoError(t, err)

	return client
}

func NewElasticacheClientE(t *testing.T, region string) (*elasticache.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return elasticache.NewFromConfig(*sess), nil
}
