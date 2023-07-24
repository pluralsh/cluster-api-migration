package aws

import (
	"context"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Migrator struct {
	accessor api.ClusterAccessor
}

func (m Migrator) AddTags(tags map[string]string) error {
	if err := m.accessor.AddClusterTags(tags); err != nil {
		return err
	}
	if err := m.accessor.AddMachinePollsTags(tags); err != nil {
		return err
	}
	if err := m.accessor.AddVirtualNetworkTags(tags); err != nil {
		return err
	}
	return nil
}

func (m Migrator) Convert() (*api.Values, error) {
	c, err := m.accessor.GetCluster()
	if err != nil {
		return nil, err
	}
	w, err := m.accessor.GetWorkers()
	if err != nil {
		return nil, err
	}
	return &api.Values{
		Provider: api.ClusterProviderAWS,
		Type:     api.ClusterTypeManaged,
		Cluster:  *c,
		Workers:  *w,
	}, nil
}

func NewAWSMigrator(configuration *api.AWSConfiguration) (api.Migrator, error) {
	a, err := (&ClusterAccessor{
		configuration: configuration,
		ctx:           context.Background(),
	}).init()
	if err != nil {
		return nil, err
	}
	return &Migrator{
		accessor: a,
	}, nil
}
