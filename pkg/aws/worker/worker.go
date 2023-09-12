package worker

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tageks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go/aws/session"
	ekssdk "github.com/aws/aws-sdk-go/service/eks"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/weaveworks/eksctl/pkg/eks"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

type Worker struct {
	configuration   *api.AWSConfiguration
	ctx             context.Context
	ClusterProvider *eks.ClusterProvider
}

func NewAWSWorker(ctx context.Context, configuration *api.AWSConfiguration, clusterProvider *eks.ClusterProvider) *Worker {
	return &Worker{
		configuration:   configuration,
		ctx:             ctx,
		ClusterProvider: clusterProvider,
	}
}

func (this *Worker) GetWorkers() (*api.Workers, error) {
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
		WorkersSpec: api.WorkersSpec{
			AWSWorkers: getDefaultAwsWorkers(),
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
		newWorkers[*ng] = &api.AWSWorker{
			Replicas:    int(*nodeGroup.Nodegroup.ScalingConfig.DesiredSize),
			Labels:      nil,
			Annotations: nil,
			IsMultiAZ:   true, // default to true so that the availability zones we discovered are used
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

func (this *Worker) AddMachinePollsTags(tags map[string]string) error {
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

func getDefaultAwsWorkers() *api.AWSWorkers {
	return &api.AWSWorkers{
		"small-burst-on-demand":  nil,
		"medium-burst-on-demand": nil,
		"large-burst-on-demand":  nil,
	}
}

func taintEffect(t string) api.TaintEffect {
	if t == "NO_SCHEDULE" {
		return api.TaintEffectNoSchedule
	} else if t == "NO_EXECUTE" {
		return api.TaintEffectNoExecute
	}
	return api.TaintEffectPreferNoSchedule
}
