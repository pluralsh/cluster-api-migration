package cluster

import "github.com/pluralsh/cluster-api-migration/pkg/api"

func (this *Cluster) VirtualNetwork() api.ManagedControlPlaneVirtualNetwork {
	// np := this.Cluster.Properties.NetworkProfile

	return api.ManagedControlPlaneVirtualNetwork{
		Name: "",
		// CIDRBlock:     nil,
		// Subnet:        nil,
		ResourceGroup: this.ResourceGroup,
	}
}
