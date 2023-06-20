package cluster

import (
	"cloud.google.com/go/container/apiv1/containerpb"
	"google.golang.org/api/compute/v1"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

func (this *Cluster) DatapathProvider() api.DatapathProvider {
	if this.NetworkConfig == nil {
		return api.DatapathProvider_UNSPECIFIED
	}

	switch this.NetworkConfig.DatapathProvider {
	case containerpb.DatapathProvider_DATAPATH_PROVIDER_UNSPECIFIED:
		return api.DatapathProvider_UNSPECIFIED
	case containerpb.DatapathProvider_LEGACY_DATAPATH:
		return api.DatapathProvider_LEGACY_DATAPATH
	case containerpb.DatapathProvider_ADVANCED_DATAPATH:
		return api.DatapathProvider_ADVANCED_DATAPATH
	}

	return api.DatapathProvider_UNSPECIFIED
}

func (this *Cluster) Network() *api.GCPNetwork {
	return &api.GCPNetwork{
		Name:                  this.GetNetwork(),
		AutoCreateSubnetworks: this.AutoCreateSubnetworks(),
		DatapathProvider:      this.DatapathProvider(),
	}
}

func (this *Cluster) AutoCreateSubnetworks() bool {
	if this.network == nil {
		return false
	}

	return this.network.AutoCreateSubnetworks
}

func (this *Cluster) Subnets() api.GCPSubnets {
	result := make([]api.GCPSubnet, 0)
	for _, subnet := range this.subnetworks {
		result = append(result, this.toSubnet(subnet))
	}

	return result
}

func (this *Cluster) toSubnet(subnetwork *compute.Subnetwork) api.GCPSubnet {
	secondaryCidrBlocks := map[string]string{}
	for _, ipRange := range subnetwork.SecondaryIpRanges {
		secondaryCidrBlocks[ipRange.RangeName] = ipRange.IpCidrRange
	}

	return api.GCPSubnet{
		Name:                subnetwork.Name,
		CidrBlock:           subnetwork.IpCidrRange,
		Description:         subnetwork.Description,
		SecondaryCidrBlocks: secondaryCidrBlocks,
		PrivateGoogleAccess: subnetwork.PrivateIpGoogleAccess,
		EnableFlowLogs:      subnetwork.EnableFlowLogs,
		Purpose:             api.GCPSubnetPurpose(subnetwork.Purpose),
	}
}

func (this *Cluster) CIDRBlocks() []string {
	// TODO: Check if this should be read from node pools
	return []string{this.ClusterIpv4Cidr}
}
