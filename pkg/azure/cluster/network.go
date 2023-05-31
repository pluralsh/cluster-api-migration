package cluster

import "github.com/pluralsh/cluster-api-migration/pkg/api"

func (this *Cluster) Network() *api.GCPNetwork {
	return &api.GCPNetwork{
		AutoCreateSubnetworks: this.AutoCreateSubnetworks(),
	}
}

func (this *Cluster) AutoCreateSubnetworks() bool {
	// TODO: Add logic
	return true
}

func (this *Cluster) Subnets() api.GCPSubnets {
	// TODO: Add logic
	return []api.GCPSubnet{}
}

func (this *Cluster) CIDRBlocks() []string {
	// TODO: Add logic
	return []string{}
}
