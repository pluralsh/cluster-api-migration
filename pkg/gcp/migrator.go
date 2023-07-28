package gcp

import (
	"context"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Migrator struct {
	accessor api.ClusterAccessor
}

func (this *Migrator) Destroy() error {
	return nil
}

func (this *Migrator) AddTags(tags map[string]string) error {
	if err := this.accessor.AddClusterTags(tags); err != nil {
		return err
	}
	if err := this.accessor.AddMachinePollsTags(tags); err != nil {
		return err
	}
	return nil
}

func (this *Migrator) Convert() (*api.Values, error) {
	c, err := this.accessor.GetCluster()
	if err != nil {
		return nil, err
	}

	w, err := this.accessor.GetWorkers()
	if err != nil {
		return nil, err
	}

	return &api.Values{
		Provider: api.ClusterProviderGoogle,
		Type:     api.ClusterTypeManaged,
		Cluster:  *c,
		Workers:  *w,
	}, nil
}

func NewGCPMigrator(configuration *api.GCPConfiguration) (api.Migrator, error) {
	a, err := (&ClusterAccessor{
		configuration: configuration,
		ctx:           context.Background(),
	}).init()

	return &Migrator{accessor: a}, err
}
