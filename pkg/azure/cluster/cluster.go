package cluster

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2022-03-01/containerservice"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
)

type Cluster struct {
	Cluster        *containerservice.ManagedCluster
	VNet           *armnetwork.VirtualNetwork
	ResourceGroup  string
	SubscriptionID string
}

func (cluster *Cluster) SKU() *api.AKSSku {
	if cluster.Cluster.Sku == nil {
		return nil
	}

	return &api.AKSSku{Tier: (*string)(&cluster.Cluster.Sku.Tier)}
}

func (cluster *Cluster) SSHPublicKey() *string {
	lp := cluster.Cluster.LinuxProfile
	if lp != nil && lp.SSH != nil && lp.SSH.PublicKeys != nil && len(*lp.SSH.PublicKeys) > 0 {
		return (*lp.SSH.PublicKeys)[0].KeyData
	}

	return resources.Ptr("skip")
}

func (cluster *Cluster) Convert() (*api.Cluster, error) {
	return &api.Cluster{
		Name:              *cluster.Cluster.Name,
		PodCIDRBlocks:     cluster.PodCIDRBlocks(),
		ServiceCIDRBlocks: cluster.ServiceCIDRBlocks(),
		KubernetesVersion: *cluster.Cluster.KubernetesVersion,
		CloudSpec: api.CloudSpec{
			AzureCloudSpec: &api.AzureCloudSpec{
				// Omitted client ID and secret as it will be filled by values.yaml.tpl.
				ClusterIdentityName:    "cluster-identity",
				ClusterIdentityType:    "ServicePrincipal",
				AllowedNamespaces:      &api.AllowedNamespaces{},
				TenantID:               *cluster.Cluster.Identity.TenantID,
				ClientSecretName:       "cluster-identity-secret",
				SubscriptionID:         cluster.SubscriptionID,
				Location:               *cluster.Cluster.Location,
				ResourceGroupName:      cluster.ResourceGroup,
				NodeResourceGroupName:  *cluster.Cluster.NodeResourceGroup,
				VirtualNetwork:         cluster.VirtualNetwork(),
				NetworkPlugin:          (*string)(&cluster.Cluster.NetworkProfile.NetworkPlugin),
				NetworkPolicy:          (*string)(&cluster.Cluster.NetworkProfile.NetworkPolicy),
				OutboundType:           (*string)(&cluster.Cluster.NetworkProfile.OutboundType),
				DNSServiceIP:           cluster.Cluster.NetworkProfile.DNSServiceIP,
				SSHPublicKey:           cluster.SSHPublicKey(),
				SKU:                    cluster.SKU(),
				LoadBalancerSKU:        (*string)(&cluster.Cluster.NetworkProfile.LoadBalancerSku),
				LoadBalancerProfile:    cluster.LoadBalancerProfile(),
				APIServerAccessProfile: cluster.APIServerAccessProfile(),
				AutoScalerProfile:      cluster.AutoscalerProfile(),
				AddonProfiles:          cluster.AddonProfiles(),
			},
		},
	}, nil
}

func NewAzureCluster(subscriptionId, resourceGroup string,
	cluster *containerservice.ManagedCluster, vnet *armnetwork.VirtualNetwork) *Cluster {
	return &Cluster{
		Cluster:        cluster,
		VNet:           vnet,
		ResourceGroup:  resourceGroup,
		SubscriptionID: subscriptionId,
	}
}
