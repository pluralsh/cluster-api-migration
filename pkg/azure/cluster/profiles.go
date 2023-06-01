package cluster

import "github.com/pluralsh/cluster-api-migration/pkg/api"

func (cluster *Cluster) AutoscalerProfile() *api.AutoScalerProfile {
	ap := cluster.Cluster.Properties.AutoScalerProfile
	if ap == nil {
		return nil
	}

	return &api.AutoScalerProfile{
		BalanceSimilarNodeGroups:      ap.BalanceSimilarNodeGroups,
		Expander:                      (*string)(ap.Expander),
		MaxEmptyBulkDelete:            ap.MaxEmptyBulkDelete,
		MaxGracefulTerminationSec:     ap.MaxGracefulTerminationSec,
		MaxNodeProvisionTime:          ap.MaxNodeProvisionTime,
		MaxTotalUnreadyPercentage:     ap.MaxTotalUnreadyPercentage,
		NewPodScaleUpDelay:            ap.NewPodScaleUpDelay,
		OkTotalUnreadyCount:           ap.OkTotalUnreadyCount,
		ScaleDownDelayAfterAdd:        ap.ScaleDownDelayAfterAdd,
		ScaleDownDelayAfterDelete:     ap.ScaleDownDelayAfterDelete,
		ScaleDownDelayAfterFailure:    ap.ScaleDownDelayAfterFailure,
		ScaleDownUnneededTime:         ap.ScaleDownUnneededTime,
		ScaleDownUnreadyTime:          ap.ScaleDownUnreadyTime,
		ScaleDownUtilizationThreshold: ap.ScaleDownUtilizationThreshold,
		ScanInterval:                  ap.ScanInterval,
		SkipNodesWithLocalStorage:     ap.SkipNodesWithLocalStorage,
		SkipNodesWithSystemPods:       ap.SkipNodesWithSystemPods,
	}
}

func (cluster *Cluster) APIServerAccessProfile() *api.APIServerAccessProfile {
	asap := cluster.Cluster.Properties.APIServerAccessProfile
	if asap == nil {
		return nil
	}

	return &api.APIServerAccessProfile{
		AuthorizedIPRanges:             asap.AuthorizedIPRanges,
		EnablePrivateCluster:           asap.EnablePrivateCluster,
		PrivateDNSZone:                 asap.PrivateDNSZone,
		EnablePrivateClusterPublicFQDN: asap.EnablePrivateClusterPublicFQDN,
	}
}

func (cluster *Cluster) LoadBalancerProfileOutboundIPPrefixes() []string {
	prefixes := []string{}
	for _, prefix := range cluster.Cluster.Properties.NetworkProfile.LoadBalancerProfile.OutboundIPPrefixes.PublicIPPrefixes {
		prefixes = append(prefixes, *prefix.ID)
	}

	return prefixes
}

func (cluster *Cluster) LoadBalancerProfileOutboundIPs() []string {
	ips := []string{}
	for _, ip := range cluster.Cluster.Properties.NetworkProfile.LoadBalancerProfile.OutboundIPs.PublicIPs {
		ips = append(ips, *ip.ID)
	}

	return ips
}

func (cluster *Cluster) LoadBalancerProfile() *api.LoadBalancerProfile {
	lbp := cluster.Cluster.Properties.NetworkProfile.LoadBalancerProfile
	if lbp == nil {
		return nil
	}

	return &api.LoadBalancerProfile{
		ManagedOutboundIPs:     lbp.ManagedOutboundIPs.Count,
		OutboundIPPrefixes:     cluster.LoadBalancerProfileOutboundIPPrefixes(),
		OutboundIPs:            cluster.LoadBalancerProfileOutboundIPs(),
		AllocatedOutboundPorts: lbp.AllocatedOutboundPorts,
		IdleTimeoutInMinutes:   lbp.IdleTimeoutInMinutes,
	}
}
