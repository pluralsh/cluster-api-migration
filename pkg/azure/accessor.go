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

func (accessor *ClusterAccessor) init() (api.ClusterAccessor, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armresources.NewClientFactory(accessor.configuration.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	accessor.resourceGroupClient = clientFactory.NewResourceGroupsClient()

	csClientFactory, err := armcontainerservice.NewClientFactory(accessor.configuration.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	accessor.managedClustersClient = csClientFactory.NewManagedClustersClient()
	accessor.agentPoolsClient = csClientFactory.NewAgentPoolsClient()

	return accessor, nil
}

func (accessor *ClusterAccessor) GetCluster() (*api.Cluster, error) {
	c, err := accessor.managedClustersClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, accessor.configuration.Name, nil)
	if err != nil {
		return nil, err
	}

	azureCluster := cluster.NewAzureCluster(accessor.configuration.SubscriptionID, accessor.configuration.ResourceGroup, &c.ManagedCluster)
	return azureCluster.Convert()
}

// TODO: Avoid connecting Azure API twice.
func (accessor *ClusterAccessor) GetWorkers() (*api.Workers, error) {
	c, err := accessor.managedClustersClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, accessor.configuration.Name, nil)
	if err != nil {
		return nil, err
	}

	azureWorkers := worker.NewAzureWorkers(accessor.configuration.SubscriptionID, accessor.configuration.ResourceGroup, &c.ManagedCluster)
	return azureWorkers.Convert()
}
