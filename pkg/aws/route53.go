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
	if hostName == "" {
		return nil, fmt.Errorf("hostName cannot be empty")
	}
	if awsRegion == "" {
		return nil, fmt.Errorf("awsRegion cannot be empty")
	}

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

	// Find exact match as ListHostedZonesByName returns zones in lexicographic order
	for _, zone := range response.HostedZones {
		return &zone, nil
	}
	return nil, fmt.Errorf("no exact match found for hosted zone %s", hostName)
}

func CleanDNSZoneID(t *testing.T, ctx context.Context, zoneID string, awsRegion string) error {
	route53Client, err := aws.NewRoute53ClientE(t, awsRegion)
	if err != nil {
		return err
	}

	o, err := route53Client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:    &zoneID,
		MaxItems:        100,
	})
	if err != nil {
		return err
	}

	var changes []types.Change

	for _, record := range o.ResourceRecordSets {
		if record.Type == types.RRTypeNs || record.Type == types.RRTypeSoa {
			continue
		}
		// Build a deletion change for each record
		changes = append(changes, types.Change{
			Action:            types.ChangeActionDelete,
			ResourceRecordSet: &record,
		})
	}

	if len(changes) == 0 {
		fmt.Println("No deletable records found.")
		return nil
	}

	// Prepare the change batch
	changeBatch := &types.ChangeBatch{
		Changes: changes,
	}

	// Call ChangeResourceRecordSets to delete the records
	changeInput := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: zoneID,
		ChangeBatch:  changeBatch,
	}

	_, err = route53Client.ChangeResourceRecordSets(ctx, changeInput)
	if err != nil {
		return err
	}

	return nil
}
