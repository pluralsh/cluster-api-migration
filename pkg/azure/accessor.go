package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2022-03-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/azure/cluster"
	"github.com/pluralsh/cluster-api-migration/pkg/azure/worker"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
)

type ClusterAccessor struct {
	configuration         *api.AzureConfiguration
	ctx                   context.Context
	managedClustersClient containerservice.ManagedClustersClient
	virtualNetworksClient *armnetwork.VirtualNetworksClient
}

func (accessor *ClusterAccessor) Destroy() error {
	return nil
}

func (accessor *ClusterAccessor) AddClusterTags(tags map[string]string) error {
	params := containerservice.TagsObject{Tags: map[string]*string{}}
	for key, value := range tags {
		params.Tags[key] = resources.Ptr(value)
	}

	_, err := accessor.managedClustersClient.UpdateTags(accessor.ctx, accessor.configuration.ResourceGroup, accessor.configuration.Name, params)
	return err
}

func (accessor *ClusterAccessor) AddMachinePollsTags(Tags map[string]string) error {
	return nil
}

func (accessor *ClusterAccessor) AddVirtualNetworkTags(tags map[string]string) error {
	c, err := accessor.managedClustersClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, accessor.configuration.Name)
	if err != nil {
		return err
	}

	vnet, _ := cluster.VirtualNetworkSubnetNames(&c)

	params := armnetwork.TagsObject{Tags: map[string]*string{}}
	for key, value := range tags {
		params.Tags[key] = resources.Ptr(value)
	}

	_, err = accessor.virtualNetworksClient.UpdateTags(accessor.ctx, accessor.configuration.ResourceGroup, vnet, params, nil)
	return err
}

func (accessor *ClusterAccessor) init() (api.ClusterAccessor, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	accessor.managedClustersClient = containerservice.NewManagedClustersClient(accessor.configuration.SubscriptionID)
	if err != nil {
		return nil, err
	}

	accessor.managedClustersClient.Authorizer, err = auth.NewAuthorizerFromCLI()
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
	c, err := accessor.managedClustersClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, accessor.configuration.Name)
	if err != nil {
		return nil, err
	}

	vnet, _ := cluster.VirtualNetworkSubnetNames(&c)
	v, err := accessor.virtualNetworksClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, vnet, nil)
	if err != nil {
		return nil, err
	}

	azureCluster := cluster.NewAzureCluster(
		accessor.configuration.SubscriptionID,
		accessor.configuration.ResourceGroup,
		&c,
		&v.VirtualNetwork)
	return azureCluster.Convert()
}

// TODO: Avoid connecting Azure API twice.
func (accessor *ClusterAccessor) GetWorkers() (*api.Workers, error) {
	c, err := accessor.managedClustersClient.Get(accessor.ctx, accessor.configuration.ResourceGroup, accessor.configuration.Name)
	if err != nil {
		return nil, err
	}

	azureWorkers := worker.NewAzureWorkers(accessor.configuration.SubscriptionID, accessor.configuration.ResourceGroup, &c)
	return azureWorkers.Convert()
}
