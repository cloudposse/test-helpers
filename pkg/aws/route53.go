package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/gruntwork-io/terratest/modules/aws"
)


func GetDNSZoneByNameE(t *testing.T, ctx context.Context, hostName string, awsRegion string) (*types.HostedZone, error) {
	client, err := aws.NewRoute53ClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	response, err := client.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{DNSName: &hostName})
	if err != nil {
		return nil, err
	}
	if len(response.HostedZones) == 0 {
		return nil, fmt.Errorf("no hosted zones found for %s", hostName)
	}
	return &response.HostedZones[0], nil
}
