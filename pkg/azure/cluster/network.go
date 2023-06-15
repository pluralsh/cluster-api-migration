package cluster

import (
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2022-03-01/containerservice"
	"strings"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

// VirtualNetworkSubnetNames reads virtual network and subnet names from agent pool profiles in form:
// /subscriptions/.../resourceGroups/.../providers/Microsoft.Network/virtualNetworks/.../subnets/...
// TODO: Find a better way to do this.
func VirtualNetworkSubnetNames(cluster *containerservice.ManagedCluster) (string, string) {
	if cluster.AgentPoolProfiles != nil {
		for _, app := range *cluster.AgentPoolProfiles {
			split := strings.Split(*app.VnetSubnetID, "/")
			if len(split) >= 10 {
				return split[8], split[10]
			}
		}
	}

	return "", ""
}

func (cluster *Cluster) PodCIDRBlocks() []string {
	cidrBlocks := make([]string, 0)
	if cluster.Cluster.NetworkProfile.PodCidrs != nil {
		for _, cidrBlock := range *cluster.Cluster.NetworkProfile.PodCidrs {
			cidrBlocks = append(cidrBlocks, cidrBlock)
		}
	}

	return cidrBlocks
}

func (cluster *Cluster) ServiceCIDRBlocks() []string {
	cidrBlocks := make([]string, 0)
	for _, cidrBlock := range *cluster.Cluster.NetworkProfile.ServiceCidrs {
		cidrBlocks = append(cidrBlocks, cidrBlock)
	}

	return cidrBlocks
}

func (cluster *Cluster) VirtualNetworkCIDRBlock() string {
	if len(cluster.VNet.Properties.AddressSpace.AddressPrefixes) > 0 {
		return *cluster.VNet.Properties.AddressSpace.AddressPrefixes[0]
	}

	return ""
}

func (cluster *Cluster) SubnetCIDRBlock() string {
	if len(cluster.VNet.Properties.Subnets) > 0 {
		return *cluster.VNet.Properties.Subnets[0].Properties.AddressPrefix
	}

	return ""
}

func (cluster *Cluster) VirtualNetwork() api.ManagedControlPlaneVirtualNetwork {
	_, subnet := VirtualNetworkSubnetNames(cluster.Cluster)

	return api.ManagedControlPlaneVirtualNetwork{
		Name:      *cluster.VNet.Name,
		CIDRBlock: cluster.VirtualNetworkCIDRBlock(),
		Subnet: api.ManagedControlPlaneSubnet{
			Name:             subnet,
			CIDRBlock:        cluster.SubnetCIDRBlock(),
			ServiceEndpoints: nil, // TODO: Do we need to fill it?
			PrivateEndpoints: nil, // TODO: Do we need to fill it?
		},
		ResourceGroup: cluster.ResourceGroup,
	}
}
