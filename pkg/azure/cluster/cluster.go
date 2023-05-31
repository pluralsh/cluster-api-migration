package cluster

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Cluster struct {
	*armcontainerservice.ManagedCluster

	ResourceGroup string
}

func (this *Cluster) KubernetesVersion() string {
	return *this.Properties.KubernetesVersion
}

func (this *Cluster) Convert() *api.Cluster {
	return &api.Cluster{
		Name:              *this.Name,
		CIDRBlocks:        nil,
		KubernetesVersion: this.KubernetesVersion(),
		CloudSpec: api.CloudSpec{
			AzureCloudSpec: &api.AzureCloudSpec{
				// TODO: Change.
				// Exported clusters will use service principal auth method,
				// not the one they were using before.
				ClusterIdentityType: "ServicePrincipal",
				ClusterIdentityName: "cluster-identity",
				ClientID:            "", // provider.clientId,
				ClientSecret:        "", // provider.clientSecret,
				ClientSecretName:    "cluster-identity-secret",
				TenantID:            *this.Identity.TenantID,
				SubscriptionID:      "", // provider.subsctiptionId,

				Location:          *this.Location,
				ResourceGroupName: this.ResourceGroup,
				SSHPublicKey:      "",
			},
		},
	}
}

func NewAzureCluster(resourceGroup string, cluster *armcontainerservice.ManagedCluster) *Cluster {
	return &Cluster{
		ManagedCluster: cluster,
		ResourceGroup:  resourceGroup,
	}
}
