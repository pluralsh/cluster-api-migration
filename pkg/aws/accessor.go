package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tageks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go/aws/session"
	ekssdk "github.com/aws/aws-sdk-go/service/eks"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

type ClusterAccessor struct {
	configuration     *api.AWSConfiguration
	ctx               context.Context
	ClusterProvider   *eks.ClusterProvider
	NodeGroupProvider *nodegroup.Manager
	AddonProvider     *addon.Manager
}

func (this *ClusterAccessor) AddClusterTags(tags map[string]string) error {
	cluster, err := this.ClusterProvider.GetCluster(this.ctx, this.configuration.ClusterName)
	if err != nil {
		return err
	}

	_, err = this.ClusterProvider.AWSProvider.EKS().TagResource(this.ctx, &tageks.TagResourceInput{
		ResourceArn: cluster.Arn,
		Tags:        tags,
	})
	if err != nil {
		return err
	}
	cfg, err := awsConfig.LoadDefaultConfig(this.ctx)
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
		Tags:      convertTags(tags),
		DryRun:    &dryFalse,
	})
	if err != nil {
		return err
	}

	subnets, err := svc.DescribeSubnets(this.ctx, &ec2.DescribeSubnetsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})

	subnetTags := map[string]string{"kubernetes.io/role/internal-elb": "1", "sigs.k8s.io/cluster-api-provider-aws/role": "common"}
	for k, v := range tags {
		subnetTags[k] = v
	}
	for _, subnet := range subnets.Subnets {
		_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
			Resources: []string{*subnet.SubnetId},
			Tags:      convertTags(subnetTags),
			DryRun:    &dryFalse,
		})
		if err != nil {
			return err
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

	sgTags := map[string]string{"sigs.k8s.io/cluster-api-provider-aws/role": "private"}
	for k, v := range tags {
		sgTags[k] = v
	}
	for _, sg := range sgroups.SecurityGroups {
		_, err = this.ClusterProvider.AWSProvider.EC2().CreateTags(this.ctx, &ec2.CreateTagsInput{
			Resources: []string{*sg.GroupId},
			Tags:      convertTags(sgTags),
			DryRun:    &dryFalse,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *ClusterAccessor) AddMachinePollsTags(tags map[string]string) error {
	mySession := session.Must(session.NewSession())
	eksSvc := ekssdk.New(mySession)
	ngList, err := eksSvc.ListNodegroups(&ekssdk.ListNodegroupsInput{
		ClusterName: &this.configuration.ClusterName,
	})
	if err != nil {
		return err
	}
	for _, ng := range ngList.Nodegroups {
		nodeGroup, err := eksSvc.DescribeNodegroup(&ekssdk.DescribeNodegroupInput{
			ClusterName:   &this.configuration.ClusterName,
			NodegroupName: ng,
		})
		if err != nil {
			return err
		}
		_, err = this.ClusterProvider.AWSProvider.EKS().TagResource(this.ctx, &tageks.TagResourceInput{
			ResourceArn: nodeGroup.Nodegroup.NodegroupArn,
			Tags:        tags,
		})
		if err != nil {
			return err
		}

	}

	return nil
}

func (this *ClusterAccessor) AddVirtualNetworkTags(tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) GetCluster() (*api.Cluster, error) {
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
			Tags:             map[string]string{fmt.Sprintf("kubernetes.io/cluster/%s", this.configuration.ClusterName): "owned"},
		}
		for _, tag := range subnet.Tags {
			sub.Tags[*tag.Key] = *tag.Value
		}

		newCluster.AWSCloudSpec.NetworkSpec.Subnets = append(newCluster.AWSCloudSpec.NetworkSpec.Subnets, sub)
	}
	return newCluster, nil
}

func taintEffect(t string) api.TaintEffect {
	if t == "NO_SCHEDULE" {
		return api.TaintEffectNoSchedule
	} else if t == "NO_EXECUTE" {
		return api.TaintEffectNoExecute
	}
	return api.TaintEffectPreferNoSchedule
}

func (this *ClusterAccessor) GetWorkers() (*api.Workers, error) {
	cfg, err := awsConfig.LoadDefaultConfig(this.ctx)
	if err != nil {
		return nil, err
	}
	cfg.Region = this.configuration.Region
	svc := ec2.NewFromConfig(cfg)
	mySession := session.Must(session.NewSession())
	eksSvc := ekssdk.New(mySession)

	ngList, err := eksSvc.ListNodegroups(&ekssdk.ListNodegroupsInput{
		ClusterName: &this.configuration.ClusterName,
	})
	if err != nil {
		return nil, err
	}

	workers := &api.Workers{
		Defaults: api.DefaultsWorker{
			AWSDefaultWorker: &api.AWSWorker{
				Replicas:    0,
				Annotations: map[string]string{"cluster.x-k8s.io/replicas-managed-by": "external-autoscaler"},
				Spec: api.AWSWorkerSpec{
					Labels:         map[string]*string{},
					AdditionalTags: map[string]string{},
					AMIType:        "AL2_x86_64",
					CapacityType:   "onDemand",
					AMIVersion:     "",
				},
			},
		},
		WorkersSpec: api.WorkersSpec{
			AWSWorkers: &api.AWSWorkers{},
		},
	}
	for _, ng := range ngList.Nodegroups {
		nodeGroup, err := eksSvc.DescribeNodegroup(&ekssdk.DescribeNodegroupInput{
			ClusterName:   &this.configuration.ClusterName,
			NodegroupName: ng,
		})
		if err != nil {
			return nil, err
		}
		availabilityZones := []string{}
		subnetID := "subnet-id"
		for _, subnet := range nodeGroup.Nodegroup.Subnets {
			subnets, err := svc.DescribeSubnets(this.ctx, &ec2.DescribeSubnetsInput{
				Filters: []ec2Types.Filter{
					{Name: &subnetID, Values: []string{*subnet}},
				},
			})
			if err != nil {
				return nil, err
			}
			availabilityZones = append(availabilityZones, *subnets.Subnets[0].AvailabilityZone)
		}
		newWorkers := *workers.AWSWorkers
		newWorkers[*ng] = api.AWSWorker{
			Replicas:    int(*nodeGroup.Nodegroup.ScalingConfig.DesiredSize),
			Labels:      nil,
			Annotations: nil,
			Spec: api.AWSWorkerSpec{
				Labels:       nodeGroup.Nodegroup.Labels,
				AMIVersion:   "", //amiVersion.Version,
				AMIType:      api.ManagedMachineAMIType(*nodeGroup.Nodegroup.AmiType),
				DiskSize:     int32(*nodeGroup.Nodegroup.DiskSize),
				InstanceType: nodeGroup.Nodegroup.InstanceTypes[0],
				Scaling: &api.ManagedMachinePoolScaling{
					MinSize: int32(*nodeGroup.Nodegroup.ScalingConfig.MinSize),
					MaxSize: int32(*nodeGroup.Nodegroup.ScalingConfig.MaxSize),
				},
				AvailabilityZones: availabilityZones,
				SubnetIDs:         nodeGroup.Nodegroup.Subnets,
				Taints: func(taints []*ekssdk.Taint) api.Taints {
					newTaints := api.Taints{}
					for _, taint := range taints {
						newTaints = append(newTaints, api.Taint{
							Effect: taintEffect(*taint.Effect),
							Key:    *taint.Key,
							Value:  *taint.Value,
						})
					}
					return newTaints
				}(nodeGroup.Nodegroup.Taints),
				UpdateConfig: nil,
				AdditionalTags: func(tags map[string]*string) infrav1.Tags {
					newTags := infrav1.Tags{fmt.Sprintf("kubernetes.io/cluster/%s", this.configuration.ClusterName): "owned"}
					for key, value := range tags {
						newTags[key] = *value
					}
					return newTags
				}(nodeGroup.Nodegroup.Tags),
			},
		}
	}
	return workers, nil
}

func (this *ClusterAccessor) init() (api.ClusterAccessor, error) {
	this.ctx = context.Background()
	cmd := &cmdutils.Cmd{}
	cfg := getCfg()
	cmd.ClusterConfig = cfg
	cmd.ClusterConfig.Metadata.Name = this.configuration.ClusterName
	cmd.ClusterConfig.Metadata.Region = this.configuration.Region
	cmd.ProviderConfig.WaitTimeout = time.Minute * 5
	clusterProvider, err := cmd.NewProviderForExistingCluster(this.ctx)
	if err != nil {
		return nil, err
	}

	if ok, err := clusterProvider.CanOperate(cfg); !ok {
		return nil, err
	}
	clientSet, err := clusterProvider.NewStdClientSet(cfg)
	if err != nil {
		return nil, err
	}
	this.NodeGroupProvider = nodegroup.New(cfg, clusterProvider, clientSet, selector.New(clusterProvider.AWSProvider.Session()))
	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)
	this.AddonProvider, err = addon.New(cmd.ClusterConfig, clusterProvider.AWSProvider.EKS(), stackManager, *cmd.ClusterConfig.IAM.WithOIDC, nil, nil)
	if err != nil {
		return nil, err
	}
	this.ClusterProvider = clusterProvider
	return this, nil
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
