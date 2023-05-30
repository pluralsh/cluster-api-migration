package gcp

import (
	"context"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Migrator struct {
	accessor api.ClusterAccessor
}

func (this *Migrator) Convert() *api.Values {
	c := this.accessor.GetCluster()
	w := this.accessor.GetWorkers()

	return &api.Values{
		Provider: api.ClusterProviderGoogle,
		// TODO: currently only managed is supported
		Type:    "managed",
		Cluster: *c,
		Workers: *w,
	}
}

func NewGCPMigrator(configuration *api.GCPConfiguration) api.Migrator {
	return &Migrator{
		accessor: (&ClusterAccessor{
			configuration: configuration,
			ctx:           context.Background(),
		}).init(),
	}
}
