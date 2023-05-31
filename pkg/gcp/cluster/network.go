package cluster

import (
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

func (this *Cluster) Network() *api.GCPNetwork {
	return &api.GCPNetwork{
		Name:                  this.GetNetwork(),
		AutoCreateSubnetworks: this.AutoCreateSubnetworks(),
	}
}

func (this *Cluster) AutoCreateSubnetworks() bool {
	if this.network == nil {
		return false
	}

	return this.network.AutoCreateSubnetworks
}

func (this *Cluster) Subnets() api.GCPSubnets {
	secondaryCidrBlocks := map[string]string{}
	for _, ipRange := range this.subnetwork.SecondaryIpRanges {
		secondaryCidrBlocks[ipRange.RangeName] = ipRange.IpCidrRange
	}

	return []api.GCPSubnet{
		{
			Name:                this.subnetwork.Name,
			CidrBlock:           this.subnetwork.IpCidrRange,
			Description:         this.subnetwork.Description,
			SecondaryCidrBlocks: secondaryCidrBlocks,
			PrivateGoogleAccess: this.subnetwork.PrivateIpGoogleAccess,
			EnableFlowLogs:      this.subnetwork.EnableFlowLogs,
			Purpose:             api.GCPSubnetPurpose(this.subnetwork.Purpose),
		},
	}
}

func (this *Cluster) CIDRBlocks() []string {
	// TODO: Check if this should be read from node pools
	return []string{this.ClusterIpv4Cidr}
}
