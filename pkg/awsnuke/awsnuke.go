// package awsnuke

// import (
// 	"fmt"
// 	"os"
// 	"sort"

// 	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
// 	tt "github.com/cloudposse/test-helpers/pkg/testing"
// 	"github.com/gruntwork-io/terratest/modules/aws"
// 	"github.com/gruntwork-io/terratest/modules/testing"
// 	"github.com/rebuy-de/aws-nuke/v2/cmd"
// 	"github.com/rebuy-de/aws-nuke/v2/pkg/awsutil"
// 	nukeconfig "github.com/rebuy-de/aws-nuke/v2/pkg/config"
// 	"github.com/rebuy-de/aws-nuke/v2/resources"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/stretchr/testify/assert"
// )

// type GenerateAwsNukeConfigConfig struct {
// 	AllRegions bool
// 	AccountId  string
// 	Regions    []string
// 	TagName    string
// 	TagValue   string
// }

// type RunAwsNukeParams struct {
// 	Config        nukeconfig.Nuke
// 	Creds         awsutil.Credentials
// 	DefaultRegion string
// 	Params        cmd.NukeParameters
// 	Verbose       bool
// }

// func getAwsNukeSupportedResources() []string {
// 	names := resources.GetListerNames()
// 	sort.Strings(names)

// 	return names
// }

// // GenerateAwsNukeConfigWithTagFilter generates an aws-nuke config object and adds a tag-based filter to every aws-nuke
// // supported resource type
// func GenerateAwsNukeConfigWithTagFilter(config GenerateAwsNukeConfigConfig) nukeconfig.Nuke {
// 	filters := nukeconfig.Filters{}
// 	for _, resource := range getAwsNukeSupportedResources() {
// 		filter := nukeconfig.Filters{
// 			resource: []nukeconfig.Filter{
// 				{
// 					Property: fmt.Sprintf("tag:%s", config.TagName),
// 					Type:     "exact",
// 					Value:    config.TagValue,
// 					Invert:   "true"},
// 			},
// 		}
// 		filters.Merge(filter)
// 	}

// 	nukeCfg := nukeconfig.Nuke{
// 		Accounts: map[string]nukeconfig.Account{
// 			config.AccountId: {
// 				Filters: filters,
// 			},
// 		},
// 		AccountBlocklist: []string{"999999999999"},
// 		Regions:          config.Regions,
// 	}

// 	return nukeCfg
// }

// // RunAwsNukeExe runs the aws-nuke as a library with the given config
// func RunAwsNukeE(params RunAwsNukeParams) error {
// 	var err error

// 	err = params.Params.Validate()
// 	if err != nil {
// 		return err
// 	}

// 	if !params.Creds.HasKeys() && !params.Creds.HasProfile() && params.DefaultRegion != "" {
// 		params.Creds.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
// 		params.Creds.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

// 		sessionToken := os.Getenv("AWS_SESSION_TOKEN")
// 		if sessionToken != "" {
// 			params.Creds.SessionToken = sessionToken
// 		}
// 	}

// 	err = params.Creds.Validate()
// 	if err != nil {
// 		return err
// 	}

// 	if params.DefaultRegion != "" {
// 		awsutil.DefaultRegionID = params.DefaultRegion
// 		switch params.DefaultRegion {
// 		case endpoints.UsEast1RegionID, endpoints.UsEast2RegionID, endpoints.UsWest1RegionID, endpoints.UsWest2RegionID:
// 			awsutil.DefaultAWSPartitionID = endpoints.AwsPartitionID
// 		case endpoints.UsGovEast1RegionID, endpoints.UsGovWest1RegionID:
// 			awsutil.DefaultAWSPartitionID = endpoints.AwsUsGovPartitionID
// 		default:
// 			if params.Config.CustomEndpoints.GetRegion(params.DefaultRegion) == nil {
// 				err = fmt.Errorf("the custom region '%s' must be specified in the configuration 'endpoints'", params.DefaultRegion)
// 				log.Error(err.Error())
// 				return err
// 			}
// 		}
// 	}

// 	account, err := awsutil.NewAccount(params.Creds, params.Config.CustomEndpoints)
// 	if err != nil {
// 		return err
// 	}

// 	n := cmd.NewNuke(params.Params, *account)

// 	n.Config = &params.Config

// 	return n.Run()
// }

// func NukeTestAccountByTag(t tt.TestingT, tagName string, tagValue string, regions []string, dryRun bool) {
// 	accountID, err := aws.GetAccountIdE(t.(testing.TestingT))
// 	assert.NoError(t, err)

// 	// run sts.getcalleridentity
// 	generateConfig := GenerateAwsNukeConfigConfig{
// 		AccountId: accountID,
// 		Regions:   regions,
// 		TagName:   tagName,
// 		TagValue:  tagValue,
// 	}
// 	config := GenerateAwsNukeConfigWithTagFilter(generateConfig)

// 	nukeParams := RunAwsNukeParams{
// 		Config:        config,
// 		Creds:         awsutil.Credentials{},
// 		DefaultRegion: "",
// 		Params: cmd.NukeParameters{
// 			ConfigPath: "GeneratedConfig",
// 			NoDryRun:   !dryRun,
// 			Force:      true,
// 			ForceSleep: 3,
// 		},
// 	}

// 	err = RunAwsNukeE(nukeParams)
// 	assert.NoError(t, err)
// }
