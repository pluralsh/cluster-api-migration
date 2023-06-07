package cluster

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

// VirtualNetworkSubnetNames reads virtual network and subnet names from agent pool profiles in form:
// /subscriptions/.../resourceGroups/.../providers/Microsoft.Network/virtualNetworks/.../subnets/...
// TODO: Find a better way to do this.
func VirtualNetworkSubnetNames(cluster *armcontainerservice.ManagedCluster) (string, string) {
	if cluster.Properties.AgentPoolProfiles != nil {
		for _, app := range cluster.Properties.AgentPoolProfiles {
			split := strings.Split(*app.VnetSubnetID, "/")
			if len(split) >= 10 {
				return split[8], split[10]
			}
		}
	}

	return "", ""
}

func (cluster *Cluster) PodCIDRBlocks() []string {
	cidrBlocks := []string{}
	for _, cidrBlock := range cluster.Cluster.Properties.NetworkProfile.PodCidrs {
		cidrBlocks = append(cidrBlocks, *cidrBlock)
	}

	return cidrBlocks
}

func (cluster *Cluster) ServiceCIDRBlocks() []string {
	cidrBlocks := []string{}
	for _, cidrBlock := range cluster.Cluster.Properties.NetworkProfile.ServiceCidrs {
		cidrBlocks = append(cidrBlocks, *cidrBlock)
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
