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
