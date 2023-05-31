package cluster

import (
	"strings"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

// TODO: Find a better way to read virtual network name.
func (cluster *Cluster) VirtualNetworkSubnetNames() (string, string) {
	if cluster.Cluster.Properties.AgentPoolProfiles != nil {
		for _, app := range cluster.Cluster.Properties.AgentPoolProfiles {
			// Form: /subscriptions/.../resourceGroups/.../providers/Microsoft.Network/virtualNetworks/.../subnets/...
			split := strings.Split(*app.VnetSubnetID, "/")
			if len(split) >= 10 {
				return split[8], split[10]
			}
		}
	}

	return "", ""
}

func (cluster *Cluster) VirtualNetwork() api.ManagedControlPlaneVirtualNetwork {
	vnet, subnet := cluster.VirtualNetworkSubnetNames()

	return api.ManagedControlPlaneVirtualNetwork{
		Name: vnet,
		// CIDRBlock:     nil,
		Subnet: api.ManagedControlPlaneSubnet{
			Name: subnet,
			// CIDRBlock:     nil,
			ServiceEndpoints: nil,
			PrivateEndpoints: nil,
		},
		ResourceGroup: cluster.ResourceGroup,
	}
}
