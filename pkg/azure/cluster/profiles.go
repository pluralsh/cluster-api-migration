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

func (cluster *Cluster) AddonProfiles() []api.AddonProfile {
	ap := cluster.Cluster.Properties.AddonProfiles
	if len(ap) < 1 {
		return nil
	}

	addonProfiles := []api.AddonProfile{}
	for key, value := range ap {
		addonProfiles = append(addonProfiles, api.AddonProfile{
			Name:    key,
			Config:  value.Config,
			Enabled: *value.Enabled,
		})
	}

	return addonProfiles
}

func (cluster *Cluster) LoadBalancerProfileOutboundIPPrefixes() []string {
	lbp := cluster.Cluster.Properties.NetworkProfile.LoadBalancerProfile
	if lbp == nil || lbp.OutboundIPPrefixes == nil {
		return nil
	}

	prefixes := []string{}
	for _, prefix := range lbp.OutboundIPPrefixes.PublicIPPrefixes {
		prefixes = append(prefixes, *prefix.ID)
	}

	return prefixes
}

func (cluster *Cluster) LoadBalancerProfileOutboundIPs() []string {
	lbp := cluster.Cluster.Properties.NetworkProfile.LoadBalancerProfile
	if lbp == nil || lbp.OutboundIPs == nil {
		return nil
	}

	ips := []string{}
	for _, ip := range lbp.OutboundIPs.PublicIPs {
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
