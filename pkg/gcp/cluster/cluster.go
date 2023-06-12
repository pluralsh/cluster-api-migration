package cluster

import (
	"cloud.google.com/go/container/apiv1/containerpb"
	"google.golang.org/api/compute/v1"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
)

type Cluster struct {
	*containerpb.Cluster

	network     *compute.Network
	subnetworks []*compute.Subnetwork
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
		PodCIDRBlocks:     this.CIDRBlocks(),
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
				AdditionalLabels:       this.additionalLabels(),
			},
		},
	}
}

func (this *Cluster) additionalLabels() *api.Labels {
	if this.ResourceLabels == nil {
		return nil
	}

	return resources.Ptr(api.Labels(this.ResourceLabels))
}

func NewGCPCluster(project string, cluster *containerpb.Cluster, network *compute.Network, subnetwork []*compute.Subnetwork) *Cluster {
	return &Cluster{
		project:     project,
		Cluster:     cluster,
		network:     network,
		subnetworks: subnetwork,
	}
}
