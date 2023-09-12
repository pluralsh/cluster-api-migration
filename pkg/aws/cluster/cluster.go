package cluster

import (
	"context"
	"fmt"
	"strings"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tageks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	"github.com/weaveworks/eksctl/pkg/eks"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

type Cluster struct {
	configuration     *api.AWSConfiguration
	ctx               context.Context
	ClusterProvider   *eks.ClusterProvider
	NodeGroupProvider *nodegroup.Manager
	AddonProvider     *addon.Manager
}

func (this *Cluster) AddClusterTags(tags map[string]string) error {
	cluster, err := this.ClusterProvider.GetCluster(this.ctx, this.configuration.ClusterName)
	if err != nil {
		return err
	}
	clusterTags := map[string]string{"sigs.k8s.io/cluster-api-provider-aws/role": "common"}
	for k, v := range tags {
		clusterTags[k] = v
	}
	_, err = this.ClusterProvider.AWSProvider.EKS().TagResource(this.ctx, &tageks.TagResourceInput{
		ResourceArn: cluster.Arn,
		Tags:        clusterTags,
	})
	if err != nil {
		return err
	}
	cfg, err := awsConfig.LoadDefaultConfig(this.ctx)
	if err != nil {
		return err
	}
	cfg.Region = this.configuration.Region
	svc := ec2.NewFromConfig(cfg)
	name := "vpc-id"
	vpcs, err := svc.DescribeVpcs(this.ctx, &ec2.DescribeVpcsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})
	if err != nil {
		return err
	}
	if len(vpcs.Vpcs) != 1 {
		return fmt.Errorf("couldn't find the VPC %s", *cluster.ResourcesVpcConfig.VpcId)
	}
	vpc := vpcs.Vpcs[0]
	dryFalse := false

	_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
		Resources: []string{*vpc.VpcId},
		Tags:      convertTags(clusterTags),
		DryRun:    &dryFalse,
	})
	if err != nil {
		return err
	}

	vpce, err := svc.DescribeVpcEndpoints(this.ctx, &ec2.DescribeVpcEndpointsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})
	if err != nil {
		return err
	}
	for _, endpoint := range vpce.VpcEndpoints {
		_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
			Resources: []string{*endpoint.VpcEndpointId},
			Tags:      convertTags(clusterTags),
			DryRun:    &dryFalse,
		})
		if err != nil {
			return err
		}
	}

	subnets, err := svc.DescribeSubnets(this.ctx, &ec2.DescribeSubnetsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})

	subnetTags := map[string]string{"kubernetes.io/role/internal-elb": "1"}
	for k, v := range tags {
		subnetTags[k] = v
	}
	for _, subnet := range subnets.Subnets {
		subnetID := "association.subnet-id"
		rt, err := svc.DescribeRouteTables(this.ctx, &ec2.DescribeRouteTablesInput{
			Filters: []ec2Types.Filter{
				{Name: &subnetID, Values: []string{*subnet.SubnetId}},
			},
		})
		if err != nil {
			return err
		}
		if len(rt.RouteTables) > 0 {
			_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
				Resources: []string{*rt.RouteTables[0].RouteTableId},
				Tags:      convertTags(tags),
				DryRun:    &dryFalse,
			})
			if err != nil {
				return err
			}
		}

		_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
			Resources: []string{*subnet.SubnetId},
			Tags:      convertTags(subnetTags),
			DryRun:    &dryFalse,
		})
		if err != nil {
			return err
		}

		subnetID = "subnet-id"
		gtws, err := svc.DescribeNatGateways(this.ctx, &ec2.DescribeNatGatewaysInput{
			Filter: []ec2Types.Filter{
				{Name: &subnetID, Values: []string{*subnet.SubnetId}},
			},
		})
		if err != nil {
			return err
		}

		for _, gtw := range gtws.NatGateways {
			_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
				Resources: []string{*gtw.NatGatewayId},
				Tags:      convertTags(tags),
				DryRun:    &dryFalse,
			})
			if err != nil {
				return err
			}
		}

	}
	sgroups, err := svc.DescribeSecurityGroups(this.ctx, &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})
	if err != nil {
		return err
	}

	for _, sg := range sgroups.SecurityGroups {
		if *sg.GroupName != "default" {
			_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
				Resources: []string{*sg.GroupId},
				Tags:      convertTags(tags),
				DryRun:    &dryFalse,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Cluster) GetCluster() (*api.Cluster, error) {
	cluster, err := this.ClusterProvider.GetCluster(this.ctx, this.configuration.ClusterName)
	if err != nil {
		return nil, err
	}

	addons, err := this.AddonProvider.GetAll(this.ctx)
	if err != nil {
		return nil, err
	}
	cfg, err := awsConfig.LoadDefaultConfig(this.ctx)
	cfg.Region = this.configuration.Region
	svc := ec2.NewFromConfig(cfg)

	name := "vpc-id"
	subnets, err := svc.DescribeSubnets(this.ctx, &ec2.DescribeSubnetsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})
	if err != nil {
		return nil, err
	}
	regionName := "region-name"
	az, err := svc.DescribeAvailabilityZones(this.ctx, &ec2.DescribeAvailabilityZonesInput{
		Filters: []ec2Types.Filter{
			{Name: &regionName, Values: []string{this.configuration.Region}},
		},
	})
	if err != nil {
		return nil, err
	}
	azLimit := len(az.AvailabilityZones)

	vpcs, err := svc.DescribeVpcs(this.ctx, &ec2.DescribeVpcsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(vpcs.Vpcs) != 1 {
		return nil, fmt.Errorf("couldn't find the VPC %s", *cluster.ResourcesVpcConfig.VpcId)
	}
	vpc := vpcs.Vpcs[0]
	newCluster := &api.Cluster{
		Name:              this.configuration.ClusterName,
		PodCIDRBlocks:     []string{*vpc.CidrBlock},
		ServiceCIDRBlocks: []string{},
		KubernetesVersion: fmt.Sprintf("v%s", *cluster.Version),
		CloudSpec: api.CloudSpec{
			AWSCloudSpec: &api.AWSCloudSpec{
				Region:             this.configuration.Region,
				SecondaryCidrBlock: "",
				EndpointAccess: api.EndpointAccess{
					Public:      true,
					PublicCIDRs: nil,
					Private:     false,
				},
				RoleAdditionalPolicies: []string{},
				EncryptionConfig: api.EncryptionConfig{
					Provider:  "",
					Resources: []string{},
				},
				AdditionalTags:             map[string]string{},
				IAMAuthenticatorConfig:     api.IAMAuthenticatorConfig{},
				OIDCIdentityProviderConfig: api.OIDCIdentityProviderConfig{},
				Logging:                    api.ControlPlaneLoggingSpec{},
				SSHKeyName:                 "default",
				Version:                    "",
				RoleName:                   "",
				ControlPlaneEndpoint: clusterv1.APIEndpoint{
					Host: *cluster.Endpoint,
					Port: 443,
				},
				Labels:                nil,
				Addons:                []api.Addon{},
				AssociateOIDCProvider: false,
				Bastion: infrav1.Bastion{
					AllowedCIDRBlocks: cluster.ResourcesVpcConfig.PublicAccessCidrs,
				},
				IdentityRef: &infrav1.AWSIdentityReference{
					Name: "default",
					Kind: infrav1.ControllerIdentityKind,
				},
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID:                         *vpc.VpcId,
						CidrBlock:                  *vpc.CidrBlock,
						Tags:                       map[string]string{fmt.Sprintf("kubernetes.io/cluster/%s", this.configuration.ClusterName): "owned"},
						AvailabilityZoneSelection:  &infrav1.AZSelectionSchemeOrdered,
						AvailabilityZoneUsageLimit: &azLimit,
					},
				},
				KubeProxy: api.KubeProxy{
					Disable: false,
				},
				VpcCni: api.VpcCni{
					Disable: false,
				},
				TokenMethod: api.EKSTokenMethodIAMAuthenticator,
			},
		},
	}
	if len(vpc.Ipv6CidrBlockAssociationSet) > 0 {
		newCluster.AWSCloudSpec.NetworkSpec.VPC.IPv6 = &infrav1.IPv6{
			CidrBlock: *vpc.Ipv6CidrBlockAssociationSet[0].Ipv6CidrBlock,
			PoolID:    *vpc.Ipv6CidrBlockAssociationSet[0].Ipv6Pool,
		}
	}

	for _, vpcTag := range vpc.Tags {
		newCluster.AWSCloudSpec.NetworkSpec.VPC.Tags[*vpcTag.Key] = *vpcTag.Value
	}
	role := strings.Split(*cluster.RoleArn, "role/")
	if len(role) == 2 {
		newCluster.AWSCloudSpec.RoleName = role[1]
	}
	if cluster.Identity != nil && cluster.Identity.Oidc != nil {
		newCluster.AWSCloudSpec.AssociateOIDCProvider = true
	}
	for _, addon := range addons {
		newCluster.AWSCloudSpec.Addons = append(newCluster.AWSCloudSpec.Addons, api.Addon{
			Name:               addon.Name,
			Version:            addon.Version,
			ConflictResolution: api.AddonResolutionOverwrite,
		})
	}
	for _, subnet := range subnets.Subnets {
		subnetID := "association.subnet-id"
		rt, err := svc.DescribeRouteTables(this.ctx, &ec2.DescribeRouteTablesInput{
			Filters: []ec2Types.Filter{
				{Name: &subnetID, Values: []string{*subnet.SubnetId}},
			},
		})
		if err != nil {
			return nil, err
		}
		var rtID *string
		if len(rt.RouteTables) > 0 {
			rtID = rt.RouteTables[0].RouteTableId
		}
		subnetID = "subnet-id"
		gtw, err := svc.DescribeNatGateways(this.ctx, &ec2.DescribeNatGatewaysInput{
			Filter: []ec2Types.Filter{
				{Name: &subnetID, Values: []string{*subnet.SubnetId}},
			},
		})
		if err != nil {
			return nil, err
		}
		var gtID *string
		if len(gtw.NatGateways) > 0 {
			gtID = gtw.NatGateways[0].NatGatewayId
		}
		sub := infrav1.SubnetSpec{
			ID:               *subnet.SubnetId,
			CidrBlock:        *subnet.CidrBlock,
			AvailabilityZone: *subnet.AvailabilityZone,
			RouteTableID:     rtID,
			NatGatewayID:     gtID,
			Tags:             map[string]string{},
		}
		for _, tag := range subnet.Tags {
			sub.Tags[*tag.Key] = *tag.Value
		}

		newCluster.AWSCloudSpec.NetworkSpec.Subnets = append(newCluster.AWSCloudSpec.NetworkSpec.Subnets, sub)
	}

	sgroups, err := svc.DescribeSecurityGroups(this.ctx, &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(sgroups.SecurityGroups) > 0 {
		newCluster.AWSCloudSpec.NetworkSpec.SecurityGroupOverrides = map[infrav1.SecurityGroupRole]string{}
	}

	return newCluster, nil
}

func NewAWSCluster(ctx context.Context, configuration *api.AWSConfiguration, clusterProvider *eks.ClusterProvider, nodeGroupProvider *nodegroup.Manager, addonProvider *addon.Manager) *Cluster {
	return &Cluster{
		configuration:     configuration,
		ctx:               ctx,
		ClusterProvider:   clusterProvider,
		NodeGroupProvider: nodeGroupProvider,
		AddonProvider:     addonProvider,
	}
}

func convertTags(tags map[string]string) []types.Tag {
	ec2Tags := []types.Tag{}
	for k, v := range tags {
		key := strings.Clone(k)
		value := strings.Clone(v)
		ec2Tags = append(ec2Tags, types.Tag{
			Key:   &key,
			Value: &value,
		})
	}
	return ec2Tags
}
