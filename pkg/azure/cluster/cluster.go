package cluster

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Cluster struct {
	Cluster        *armcontainerservice.ManagedCluster
	ResourceGroup  string
	SubscriptionID string
}

func (this *Cluster) Convert() (*api.Cluster, error) {
	return &api.Cluster{
		Name:              *this.Cluster.Name,
		CIDRBlocks:        nil,
		KubernetesVersion: *this.Cluster.Properties.KubernetesVersion,
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

				TenantID:               *this.Cluster.Identity.TenantID,
				SubscriptionID:         this.SubscriptionID,
				Location:               *this.Cluster.Location,
				ResourceGroupName:      this.ResourceGroup,
				NodeResourceGroupName:  *this.Cluster.Properties.NodeResourceGroup,
				VirtualNetwork:         this.VirtualNetwork(),
				NetworkPlugin:          (*string)(this.Cluster.Properties.NetworkProfile.NetworkPlugin),
				NetworkPolicy:          (*string)(this.Cluster.Properties.NetworkProfile.NetworkPolicy),
				OutboundType:           nil,
				DNSServiceIP:           this.Cluster.Properties.NetworkProfile.DNSServiceIP,
				SSHPublicKey:           "",
				SKU:                    nil,
				LoadBalancerSKU:        (*string)(this.Cluster.Properties.NetworkProfile.LoadBalancerSKU),
				LoadBalancerProfile:    nil,
				APIServerAccessProfile: this.APIServerAccessProfile(),
				AutoScalerProfile:      this.AutoscalerProfile(),
				AADProfile:             nil,
				AddonProfiles:          nil,
			},
		},
	}, nil
}

func NewAzureCluster(subscriptionId, resourceGroup string, cluster *armcontainerservice.ManagedCluster) *Cluster {
	return &Cluster{
		Cluster:        cluster,
		ResourceGroup:  resourceGroup,
		SubscriptionID: subscriptionId,
	}
}
