package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)


func NeWAFClient(t *testing.T, region string) *wafv2.Client {
	client, err := NeWAFClientE(t, region)
	require.NoError(t, err)

	return client
}

func NeWAFClientE(t *testing.T, region string) (*wafv2.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return wafv2.NewFromConfig(*sess), nil
}

func CreateRegexPatternSet(client *wafv2.Client, name string, description string, patterns []string) (*types.RegexPatternSetSummary, error) {
	regexSetInput := &wafv2.CreateRegexPatternSetInput{
		Name:        &name,
		Scope:       types.ScopeRegional,
		Description: &description,
		RegularExpressionList: []types.Regex{
			{RegexString: &patterns[0]},
			{RegexString: &patterns[1]},
		},
	}

	output, err := client.CreateRegexPatternSet(context.Background(), regexSetInput)
	if err != nil {
		return nil, err
	}

	return output.Summary, nil
}

func GetIPSetByARN(t *testing.T, client *wafv2.Client, arn string) *wafv2.GetIPSetOutput {
	ipSets, err := client.ListIPSets(context.Background(), &wafv2.ListIPSetsInput{
		Scope: types.ScopeRegional,
	})

	var ipSetId string
	var ipSetName string
	for _, ipSet := range ipSets.IPSets {
		if *ipSet.ARN == arn {
			ipSetId = *ipSet.Id
			ipSetName = *ipSet.Name
			break
		}
	}

	ipSet, err := client.GetIPSet(context.Background(), &wafv2.GetIPSetInput{
		Id:    &ipSetId,
		Name:  &ipSetName,
		Scope: types.ScopeRegional,
	})
	require.NoError(t, err)
	require.NotNil(t, ipSet)
	require.NotNil(t, ipSet.IPSet)
	return ipSet
}

