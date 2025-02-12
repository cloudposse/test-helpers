package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/aws/aws-sdk-go-v2/service/amplify/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)

func NewAmplifyClient(t *testing.T, region string) *amplify.Client {
	client, err := NewAmplifyClientE(t, region)
	require.NoError(t, err)

	return client
}

func NewAmplifyClientE(t *testing.T, region string) (*amplify.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return amplify.NewFromConfig(*sess), nil
}

func StartDeploymentJob(t *testing.T, client *amplify.Client, id *string, branchName *string) *string {
	branch, err := client.GetBranch(context.Background(), &amplify.GetBranchInput{
		AppId:      id,
		BranchName: branchName,
	})
	require.NoError(t, err)

	var jobType types.JobType
	if branch.Branch.ActiveJobId == nil {
		jobType = types.JobTypeRelease
	} else {
		jobType = types.JobTypeRetry
	}
	jobStart, err := client.StartJob(context.Background(), &amplify.StartJobInput{
		AppId:      id,
		BranchName: branchName,
		JobId:      branch.Branch.ActiveJobId,
		JobType:    jobType,
	})
	require.NoError(t, err)
	return jobStart.JobSummary.JobId
}

