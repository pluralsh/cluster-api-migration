package aws

import (
	"context"
	"time"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/aws/cluster"
	"github.com/pluralsh/cluster-api-migration/pkg/aws/worker"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type ClusterAccessor struct {
	configuration *api.AWSConfiguration
	cluster       *cluster.Cluster
	worker        *worker.Worker
}

func (this *ClusterAccessor) AddClusterTags(tags map[string]string) error {
	return this.cluster.AddClusterTags(tags)
}

func (this *ClusterAccessor) AddMachinePollsTags(tags map[string]string) error {
	return this.worker.AddMachinePollsTags(tags)
}

func (this *ClusterAccessor) AddVirtualNetworkTags(tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) GetCluster() (*api.Cluster, error) {
	return this.cluster.GetCluster()
}

func (this *ClusterAccessor) GetWorkers() (*api.Workers, error) {
	return this.worker.GetWorkers()
}

func (this *ClusterAccessor) init() (api.ClusterAccessor, error) {
	ctx := context.Background()
	cmd := &cmdutils.Cmd{}
	cfg := getCfg()
	cmd.ClusterConfig = cfg
	cmd.ClusterConfig.Metadata.Name = this.configuration.ClusterName
	cmd.ClusterConfig.Metadata.Region = this.configuration.Region
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
	nodeGroupProvider := nodegroup.New(cfg, clusterProvider, clientSet, selector.New(clusterProvider.AWSProvider.Session()))
	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)
	addonProvider, err := addon.New(cmd.ClusterConfig, clusterProvider.AWSProvider.EKS(), stackManager, *cmd.ClusterConfig.IAM.WithOIDC, nil, nil)
	if err != nil {
		return nil, err
	}

	this.cluster = cluster.NewAWSCluster(ctx, this.configuration, clusterProvider, nodeGroupProvider, addonProvider)
	this.worker = worker.NewAWSWorker(ctx, this.configuration, clusterProvider)
	return this, nil
}
