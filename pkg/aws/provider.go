package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws/session"
	ekssdk "github.com/aws/aws-sdk-go/service/eks"
	clusterapi "github.com/pluralsh/cluster-api-migration/pkg/api"
	migrationapi "github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

type Provider struct {
	ClusterProvider   *eks.ClusterProvider
	NodeGroupProvider *nodegroup.Manager
	AddonProvider     *addon.Manager
}

func getCfg() *api.ClusterConfig {
	cfg := api.NewClusterConfig()
	cfg.IAM.WithOIDC = api.Enabled()
	return cfg
}

func GetProvider(ctx context.Context, clusterName, region string) (*Provider, error) {
	cmd := &cmdutils.Cmd{}
	cfg := getCfg()
	cmd.ClusterConfig = cfg
	cmd.ClusterConfig.Metadata.Name = clusterName
	cmd.ClusterConfig.Metadata.Region = region
	cmd.ProviderConfig.WaitTimeout = time.Minute * 5
	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
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
	nodeGroupManager := nodegroup.New(cfg, clusterProvider, clientSet, selector.New(clusterProvider.AWSProvider.Session()))
	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)
	addonManager, err := addon.New(cmd.ClusterConfig, clusterProvider.AWSProvider.EKS(), stackManager, *cmd.ClusterConfig.IAM.WithOIDC, nil, nil)
	return &Provider{
		ClusterProvider:   clusterProvider,
		NodeGroupProvider: nodeGroupManager,
		AddonProvider:     addonManager,
	}, nil
}

func GetCluster(ctx context.Context, clusterName, region string) (*migrationapi.ClusterAPI, error) {
	provider, err := GetProvider(ctx, clusterName, region)
	if err != nil {
		return nil, err
	}

	cluster, err := provider.ClusterProvider.GetCluster(ctx, clusterName)
	if err != nil {
		return nil, err
	}
	addons, err := provider.AddonProvider.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	cfg.Region = region
	svc := ec2.NewFromConfig(cfg)
	mySession := session.Must(session.NewSession())
	eksSvc := ekssdk.New(mySession)
	if err != nil {
		return nil, err
	}
	ngList, err := eksSvc.ListNodegroups(&ekssdk.ListNodegroupsInput{
		ClusterName: &clusterName,
	})
	if err != nil {
		return nil, err
	}

	name := "vpc-id"
	subnets, err := svc.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: []ec2Types.Filter{
			{Name: &name, Values: []string{*cluster.ResourcesVpcConfig.VpcId}},
		},
	})
	if err != nil {
		return nil, err
	}
	regionName := "region-name"
	az, err := svc.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{
		Filters: []ec2Types.Filter{
			{Name: &regionName, Values: []string{region}},
		},
	})
	if err != nil {
		return nil, err
	}
	azLimit := len(az.AvailabilityZones)

	vpcs, err := svc.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
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
	newCluster := &migrationapi.ClusterAPI{
		Provider: "aws",
		Type:     "managed",
		Cluster: migrationapi.Cluster{
			Name:              clusterName,
			CIDRBlocks:        []string{*vpc.CidrBlock},
			KubernetesVersion: fmt.Sprintf("v%s", *cluster.Version),
			CloudSpec: migrationapi.CloudSpec{
				AWSCloudSpec: &clusterapi.AWSCloudSpec{
					Region:     region,
					SSHKeyName: "default",
					Version:    "",
					RoleName:   "",
					ControlPlaneEndpoint: clusterv1.APIEndpoint{
						Host: *cluster.Endpoint,
						Port: 443,
					},
					Addons:                []clusterapi.Addon{},
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
							Tags:                       map[string]string{},
							AvailabilityZoneSelection:  &infrav1.AZSelectionSchemeOrdered,
							AvailabilityZoneUsageLimit: &azLimit,
						},
					},
					KubeProxy: migrationapi.KubeProxy{
						Disable: false,
					},
					VpcCni: migrationapi.VpcCni{
						Disable: false,
					},
					TokenMethod: migrationapi.EKSTokenMethodIAMAuthenticator,
				},
			},
		},
		Workers: migrationapi.Workers{
			Defaults: migrationapi.DefaultsWorker{
				AWSDefaultWorker: &migrationapi.AWSDefaultWorker{
					migrationapi.AWSWorker{
						Replicas:    0,
						Annotations: map[string]string{"cluster.x-k8s.io/replicas-managed-by": "external-autoscaler"},
						Spec: migrationapi.AWSWorkerSpec{
							AMIType:      "AL2_x86_64",
							CapacityType: "onDemand",
						},
					},
				},
			},
			WorkersSpec: migrationapi.WorkersSpec{
				AWSWorkers: &migrationapi.AWSWorkers{},
			},
		},
	}

	for _, ng := range ngList.Nodegroups {
		nodeGroup, err := eksSvc.DescribeNodegroup(&ekssdk.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: ng,
		})
		if err != nil {
			return nil, err
		}
		availabilityZones := []string{}
		subnetID := "subnet-id"
		for _, subnet := range nodeGroup.Nodegroup.Subnets {
			subnets, err := svc.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
				Filters: []ec2Types.Filter{
					{Name: &subnetID, Values: []string{*subnet}},
				},
			})
			if err != nil {
				return nil, err
			}
			availabilityZones = append(availabilityZones, *subnets.Subnets[0].AvailabilityZone)
		}
		workers := *newCluster.Workers.AWSWorkers
		workers[*ng] = migrationapi.AWSWorker{
			Replicas:    int(*nodeGroup.Nodegroup.ScalingConfig.DesiredSize),
			Labels:      nil,
			Annotations: nil,
			Spec: migrationapi.AWSWorkerSpec{
				Labels:       nodeGroup.Nodegroup.Labels,
				AMIVersion:   *nodeGroup.Nodegroup.Version,
				AMIType:      migrationapi.ManagedMachineAMIType(*nodeGroup.Nodegroup.AmiType),
				DiskSize:     int32(*nodeGroup.Nodegroup.DiskSize),
				InstanceType: nodeGroup.Nodegroup.InstanceTypes[0],
				Scaling: &migrationapi.ManagedMachinePoolScaling{
					MinSize: int32(*nodeGroup.Nodegroup.ScalingConfig.MinSize),
					MaxSize: int32(*nodeGroup.Nodegroup.ScalingConfig.MaxSize),
				},
				AvailabilityZones: availabilityZones,
				SubnetIDs:         nodeGroup.Nodegroup.Subnets,
				Taints:            nil,
				UpdateConfig:      nil,
				AdditionalTags: func(tags map[string]*string) infrav1.Tags {
					newTags := infrav1.Tags{}
					for key, value := range tags {
						newTags[key] = *value
					}
					return newTags
				}(nodeGroup.Nodegroup.Tags),
			},
		}
	}

	for _, vpcTag := range vpc.Tags {
		newCluster.Cluster.AWSCloudSpec.NetworkSpec.VPC.Tags[*vpcTag.Key] = *vpcTag.Value
	}

	role := strings.Split(*cluster.RoleArn, "role/")
	if len(role) == 2 {
		newCluster.Cluster.AWSCloudSpec.RoleName = role[1]
	}
	if cluster.Identity != nil && cluster.Identity.Oidc != nil {
		newCluster.Cluster.AWSCloudSpec.AssociateOIDCProvider = true
	}
	for _, addon := range addons {
		newCluster.Cluster.AWSCloudSpec.Addons = append(newCluster.Cluster.AWSCloudSpec.Addons, migrationapi.Addon{
			Name:               addon.Name,
			Version:            addon.Version,
			ConflictResolution: migrationapi.AddonResolutionOverwrite,
		})
	}
	for _, subnet := range subnets.Subnets {
		subnetID := "association.subnet-id"
		rt, err := svc.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{
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
		gtw, err := svc.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{
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

		newCluster.Cluster.AWSCloudSpec.NetworkSpec.Subnets = append(newCluster.Cluster.AWSCloudSpec.NetworkSpec.Subnets, sub)
	}

	return newCluster, nil
}
