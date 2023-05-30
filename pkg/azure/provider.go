package azure

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

const (
	providerName = "azure"
	clusterType  = "managed"
)

type Provider struct {
	subsctiptionId        string
	clientId              string
	clientSecret          string
	resourceGroupClient   *armresources.ResourceGroupsClient
	managedClustersClient *armcontainerservice.ManagedClustersClient
}

func GetProvider(ctx context.Context, clusterName, resourceGroupName string) (*Provider, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(subscriptionId) == 0 {
		log.Fatal("AZURE_SUBSCRIPTION_ID is not set")
	}

	clientId := os.Getenv("AZURE_CLIENT_ID")
	if len(clientId) == 0 {
		log.Fatal("AZURE_CLIENT_ID is not set")
	}

	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if len(clientSecret) == 0 {
		log.Fatal("AZURE_CLIENT_SECRET is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("authentication failure: %+v", err)
	}

	clientFactory, err := armresources.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("cannot create client factory: %+v", err)
	}

	resourceGroupClient := clientFactory.NewResourceGroupsClient()

	csClientFactory, err := armcontainerservice.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	managedClustersClient := csClientFactory.NewManagedClustersClient()

	return &Provider{
		subsctiptionId:        subscriptionId,
		clientId:              clientId,
		clientSecret:          clientSecret,
		resourceGroupClient:   resourceGroupClient,
		managedClustersClient: managedClustersClient,
	}, nil
}

func GetCluster(ctx context.Context, clusterName, resourceGroupName string) (*api.ClusterAPI, error) {
	provider, err := GetProvider(ctx, clusterName, resourceGroupName)
	if err != nil {
		return nil, err
	}

	cluster, err := provider.managedClustersClient.Get(ctx, resourceGroupName, clusterName, nil)
	if err != nil {
		return nil, err
	}

	outputCluster := api.Cluster{
		Name:              *cluster.Name,
		CIDRBlocks:        nil,
		KubernetesVersion: *cluster.Properties.KubernetesVersion,
		CloudSpec: api.CloudSpec{
			AzureCloudSpec: &api.AzureCloudSpec{
				// Exported clusters will use service principal auth method,
				// not the one they were using before.
				ClusterIdentityType: "ServicePrincipal",
				ClusterIdentityName: "cluster-identity",
				ClientID:            provider.clientId,
				ClientSecret:        provider.clientSecret,
				ClientSecretName:    "cluster-identity-secret",
				TenantID:            *cluster.Identity.TenantID,
				SubscriptionID:      provider.subsctiptionId,

				Location:          *cluster.Location,
				ResourceGroupName: resourceGroupName,
				SSHPublicKey:      "",
			},
		},
	}

	return &api.ClusterAPI{ // TODO: Fill.
		Provider: providerName,
		Type:     clusterType,
		Cluster:  outputCluster,
		Workers:  api.Workers{},
	}, nil
}
