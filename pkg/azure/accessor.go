package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/azure/cluster"
	"github.com/pluralsh/cluster-api-migration/pkg/azure/worker"
)

type ClusterAccessor struct {
	configuration         *api.AzureConfiguration
	ctx                   context.Context
	resourceGroupClient   *armresources.ResourceGroupsClient // TODO: Remove if it is not needed.
	managedClustersClient *armcontainerservice.ManagedClustersClient
	agentPoolsClient      *armcontainerservice.AgentPoolsClient
}

func (this *ClusterAccessor) init() (api.ClusterAccessor, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armresources.NewClientFactory(this.configuration.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	this.resourceGroupClient = clientFactory.NewResourceGroupsClient()

	csClientFactory, err := armcontainerservice.NewClientFactory(this.configuration.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	this.managedClustersClient = csClientFactory.NewManagedClustersClient()
	this.agentPoolsClient = csClientFactory.NewAgentPoolsClient()

	return this, nil
}

func (this *ClusterAccessor) GetCluster() (*api.Cluster, error) {
	c, err := this.managedClustersClient.Get(this.ctx, this.configuration.ResourceGroup, this.configuration.Name, nil)
	if err != nil {
		return nil, err
	}

	azureCluster := cluster.NewAzureCluster(this.configuration.SubscriptionID, this.configuration.ResourceGroup, &c.ManagedCluster)
	return azureCluster.Convert()
}

func (this *ClusterAccessor) GetWorkers() (*api.Workers, error) {
	return &api.Workers{
		Defaults: api.DefaultsWorker{
			AzureDefaultWorker: worker.AzureWorkerDefaults(),
		},
		WorkersSpec: api.WorkersSpec{},
	}, nil
}
