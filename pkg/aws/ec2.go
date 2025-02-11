func GetNatsByVpcIdE(t *testing.T, vpcId string, awsRegion string) ([]types.NatGateway, error) {
	client, err := aws.NewEc2ClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	filter := &types.Filter{Name: awstypes.String("vpc-id"), Values: []string{vpcId}}
	response, err := client.DescribeNatGateways(context.Background(), &ec2.DescribeNatGatewaysInput{
		Filter: []types.Filter{*filter},
	})
	if err != nil {
		return nil, err
	}
	return response.NatGateways, nil
}
