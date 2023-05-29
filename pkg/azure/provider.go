package azure

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
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

func GetCluster(ctx context.Context, clusterName, resourceGroupName string) (*struct{}, error) {
	provider, err := GetProvider(ctx, clusterName, resourceGroupName)
	if err != nil {
		return nil, err
	}

	cluster, err := provider.managedClustersClient.Get(ctx, resourceGroupName, clusterName, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(*cluster.Name)
	fmt.Println(*cluster.SKU.Name)

	return nil, nil
}
