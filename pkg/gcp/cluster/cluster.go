package cluster

import (
	"cloud.google.com/go/compute/apiv1/computepb"
	"cloud.google.com/go/container/apiv1/containerpb"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
)

type Cluster struct {
	*containerpb.Cluster

	network     *computepb.Network
	subnetworks *computepb.SubnetworkList
	project     string
}

func (this *Cluster) AutopilotEnabled() bool {
	if this.GetAutopilot() == nil {
		return false
	}

	return this.GetAutopilot().Enabled
}

func (this *Cluster) WorkloadIdentityEnabled() bool {
	if this.GetWorkloadIdentityConfig() == nil {
		return false
	}

	return len(this.GetWorkloadIdentityConfig().GetWorkloadPool()) > 0
}

func (this *Cluster) ReleaseChannel() *api.GCPReleaseChannel {
	if this.GetReleaseChannel() == nil {
		return nil
	}

	switch this.GetReleaseChannel().Channel {
	case containerpb.ReleaseChannel_UNSPECIFIED:
		return nil
	case containerpb.ReleaseChannel_RAPID:
		return resources.Ptr(api.ReleaseChannelRapid)
	case containerpb.ReleaseChannel_REGULAR:
		return resources.Ptr(api.ReleaseChannelRegular)
	case containerpb.ReleaseChannel_STABLE:
		return resources.Ptr(api.ReleaseChannelStable)
	default:
		return nil
	}
}

func (this *Cluster) KubernetesVersion() string {
	return this.GetCurrentMasterVersion()
}

func (this *Cluster) Convert() *api.Cluster {
	return &api.Cluster{
		Name:              this.GetName(),
		CIDRBlocks:        this.CIDRBlocks(),
		KubernetesVersion: this.KubernetesVersion(),
		CloudSpec: api.CloudSpec{
			GCPCloudSpec: &api.GCPCloudSpec{
				Project:                this.project,
				Region:                 this.Location,
				EnableAutopilot:        this.AutopilotEnabled(),
				EnableWorkloadIdentity: this.WorkloadIdentityEnabled(),
				ReleaseChannel:         this.ReleaseChannel(),
				Network:                this.Network(),
				Subnets:                this.Subnets(),
			},
		},
	}
}

func NewGCPCluster(project string, cluster *containerpb.Cluster) *Cluster {
	return &Cluster{
		project: project,
		Cluster: cluster,
	}
}
