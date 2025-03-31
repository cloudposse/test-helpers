package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/gruntwork-io/terratest/modules/testing"
	terratestAWS "github.com/gruntwork-io/terratest/modules/aws"	
	"github.com/stretchr/testify/require"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

// GetEksClusterE fetches information about an EKS cluster.
func GetEksClusterE(t testing.TestingT, ctx context.Context, region string, name string) (*types.Cluster, error) {
	client, err := NewEksClientE(t, region)
	if err != nil {
		return nil, err
	}
	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}
	output, err := client.DescribeCluster(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.Cluster, nil
}

// GetEksCluster fetches information about an EKS cluster.
func GetEksCluster(t testing.TestingT, ctx context.Context, region string, name string) *types.Cluster {
	cluster, err := GetEksClusterE(t, ctx, region, name)
	require.NoError(t, err)
	return cluster
}

// CreateEksClusterE creates EKS cluster in the given region under the given name.
func CreateEksClusterE(t testing.TestingT, ctx context.Context, region string, name string, roleArn string, enablePrivateAccess bool, enablePublicAccess bool, subnets []string, publicAccessCidrs []string, securityGroupIds []string) (*types.Cluster, error) {
	client := NewEksClient(t, region)
	cluster, err := client.CreateCluster(ctx, &eks.CreateClusterInput{
		Name: aws.String(name),
		ResourcesVpcConfig: &types.VpcConfigRequest{
			EndpointPublicAccess:  aws.Bool(enablePublicAccess),
			EndpointPrivateAccess: aws.Bool(enablePrivateAccess),
			SubnetIds:             subnets,
			SecurityGroupIds:      securityGroupIds,
			PublicAccessCidrs:     publicAccessCidrs,
		},
		RoleArn: aws.String(roleArn),
	})

	if err != nil {
		return nil, err
	}
	return cluster.Cluster, nil
}

// CreateEksCluster creates EKS cluster in the given region under the given name.
func CreateEksCluster(t testing.TestingT, ctx context.Context, region string, name string, roleArn string, enablePrivateAccess bool, enablePublicAccess bool, subnets []string, publicAccessCidrs []string, securityGroupIds []string) *types.Cluster {
	cluster, err := CreateEksClusterE(t, ctx, region, name, roleArn, enablePrivateAccess, enablePublicAccess, subnets, publicAccessCidrs, securityGroupIds)
	require.NoError(t, err)
	return cluster
}

// DeleteEksClusterE deletes existing EKS cluster in the given region.
func DeleteEksClusterE(t testing.TestingT, ctx context.Context, region string, cluster *types.Cluster) error {
	client := NewEksClient(t, region)
	_, err := client.DeleteCluster(ctx, &eks.DeleteClusterInput{
		Name: aws.String(*cluster.Name),
	})
	return err
}

// DeleteEksCluster deletes existing EKS cluster in the given region.
func DeleteEksCluster(t testing.TestingT, ctx context.Context, region string, cluster *types.Cluster) {
	err := DeleteEksClusterE(t, ctx, region, cluster)
	require.NoError(t, err)
}

// NewEksClient creates an EKS client.
func NewEksClient(t testing.TestingT, region string) *eks.Client {
	client, err := NewEksClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewEksClientE creates an EKS client.
func NewEksClientE(t testing.TestingT, region string) (*eks.Client, error) {
	sess, err := terratestAWS.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return eks.NewFromConfig(*sess), nil
}


func NewK8SClientConfig(cluster *types.Cluster) (*rest.Config, error) {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}
	opts := &token.GetTokenOptions{
		ClusterID: *cluster.Name,
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, err
	}

	return &rest.Config{
		Host:        *cluster.Endpoint,
		BearerToken: tok.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: ca,
		},
	}, nil
}


func NewK8SClientset(cluster *types.Cluster) (*kubernetes.Clientset, error) {
	config, err := NewK8SClientConfig(cluster)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
