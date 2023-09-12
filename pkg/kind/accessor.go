package kind

import (
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type ClusterAccessor struct {
	configuration *api.KindConfiguration
}

func (this *ClusterAccessor) AddClusterTags(tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) AddMachinePollsTags(tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) AddVirtualNetworkTags(tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) GetCluster() (*api.Cluster, error) {
	return nil, nil
}

func (this *ClusterAccessor) GetWorkers() (*api.Workers, error) {
	return nil, nil
}

func (this *ClusterAccessor) init() (api.ClusterAccessor, error) {
	return this, nil
}
