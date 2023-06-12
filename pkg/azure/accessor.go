package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/azure/cluster"
	"github.com/pluralsh/cluster-api-migration/pkg/azure/worker"
)

type ClusterAccessor struct {
	configuration         *api.AzureConfiguration
	ctx                   context.Context
	managedClustersClient *armcontainerservice.ManagedClustersClient
	virtualNetworksClient *armnetwork.VirtualNetworksClient
}

func (accessor *ClusterAccessor) init() (api.ClusterAccessor, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	accessor.managedClustersClient, err = armcontainerservice.NewManagedClustersClient(accessor.configuration.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	accessor.virtualNetworksClient, err = armnetwork.NewVirtualNetworksClient(accessor.configuration.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return accessor, nil
}

func (accessor *ClusterAccessor) GetCluster() (*api.Cluster, error) {
	c, err := accessor.managedClustersClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, accessor.configuration.Name, nil)
	if err != nil {
		return nil, err
	}

	vnet, _ := cluster.VirtualNetworkSubnetNames(&c.ManagedCluster)
	v, err := accessor.virtualNetworksClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, vnet, nil)
	if err != nil {
		return nil, err
	}

	azureCluster := cluster.NewAzureCluster(
		accessor.configuration.SubscriptionID,
		accessor.configuration.ResourceGroup,
		accessor.configuration.SSHPublicKey,
		&c.ManagedCluster,
		&v.VirtualNetwork)
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
