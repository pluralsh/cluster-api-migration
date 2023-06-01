package cluster

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Cluster struct {
	Cluster        *armcontainerservice.ManagedCluster
	VNet           *armnetwork.VirtualNetwork
	ResourceGroup  string
	SubscriptionID string
}

func (cluster *Cluster) SKU() *api.AKSSku {
	if cluster.Cluster.SKU == nil {
		return nil
	}

	return &api.AKSSku{Tier: (*string)(cluster.Cluster.SKU.Tier)}
}

func (cluster *Cluster) Convert() (*api.Cluster, error) {
	return &api.Cluster{
		Name:              *cluster.Cluster.Name,
		PodCIDRBlocks:     cluster.PodCIDRBlocks(),
		ServiceCIDRBlocks: cluster.ServiceCIDRBlocks(),
		KubernetesVersion: *cluster.Cluster.Properties.KubernetesVersion,
		CloudSpec: api.CloudSpec{
			AzureCloudSpec: &api.AzureCloudSpec{
				// TODO: Change.
				// Exported clusters will use service principal auth method,
				// not the one they were using before.
				ClusterIdentityType: "ServicePrincipal",
				ClusterIdentityName: "cluster-identity",
				AllowedNamespaces:   nil,
				ClientID:            "", // provider.clientId,
				ClientSecret:        "", // provider.clientSecret,
				ClientSecretName:    "cluster-identity-secret",
				ResourceID:          "",

				TenantID:               *cluster.Cluster.Identity.TenantID,
				SubscriptionID:         cluster.SubscriptionID,
				Location:               *cluster.Cluster.Location,
				ResourceGroupName:      cluster.ResourceGroup,
				NodeResourceGroupName:  *cluster.Cluster.Properties.NodeResourceGroup,
				VirtualNetwork:         cluster.VirtualNetwork(),
				NetworkPlugin:          (*string)(cluster.Cluster.Properties.NetworkProfile.NetworkPlugin),
				NetworkPolicy:          (*string)(cluster.Cluster.Properties.NetworkProfile.NetworkPolicy),
				OutboundType:           (*string)(cluster.Cluster.Properties.NetworkProfile.OutboundType),
				DNSServiceIP:           cluster.Cluster.Properties.NetworkProfile.DNSServiceIP,
				SSHPublicKey:           "",
				SKU:                    cluster.SKU(),
				LoadBalancerSKU:        (*string)(cluster.Cluster.Properties.NetworkProfile.LoadBalancerSKU),
				LoadBalancerProfile:    nil,
				APIServerAccessProfile: cluster.APIServerAccessProfile(),
				AutoScalerProfile:      cluster.AutoscalerProfile(),
				AADProfile:             nil,
				AddonProfiles:          nil,
			},
		},
	}, nil
}

func NewAzureCluster(subscriptionId, resourceGroup string, cluster *armcontainerservice.ManagedCluster, vnet *armnetwork.VirtualNetwork) *Cluster {
	return &Cluster{
		Cluster:        cluster,
		VNet:           vnet,
		ResourceGroup:  resourceGroup,
		SubscriptionID: subscriptionId,
	}
}
