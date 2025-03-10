package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)


func NewWAFClient(t *testing.T, region string) *wafv2.Client {
	client, err := NewWAFClientE(t, region)
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
	if len(patterns) < 1 {
		return nil, fmt.Errorf("at least one pattern is required")
	}

	var regexList []types.Regex
	for _, pattern := range patterns {
		p := pattern // Create a new variable to avoid loop variable capture issues
		regexList = append(regexList, types.Regex{RegexString: &p})
	}

	regexSetInput := &wafv2.CreateRegexPatternSetInput{
		Name:                  &name,
		Scope:                 types.ScopeRegional,
		Description:           &description,
		RegularExpressionList: regexList,
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
	require.NoError(t, err)
	require.NotNil(t, ipSets)

	var ipSetId string
	var ipSetName string
	found := false
	for _, ipSet := range ipSets.IPSets {
		if *ipSet.ARN == arn {
			ipSetId = *ipSet.Id
			ipSetName = *ipSet.Name
			found = true
			break
		}
	}
	require.True(t, found, "IP set with ARN %s not found", arn)

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

