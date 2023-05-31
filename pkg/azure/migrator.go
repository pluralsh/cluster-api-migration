package azure

import (
	"context"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Migrator struct {
	accessor api.ClusterAccessor
}

func (migrator *Migrator) Convert() *api.Values {
	c := migrator.accessor.GetCluster()
	w := migrator.accessor.GetWorkers()

	return &api.Values{
		Provider: api.ClusterProviderAzure,
		Type:     api.ClusterTypeManaged,
		Cluster:  *c,
		Workers:  *w,
	}
}

func NewAzureMigrator(configuration *api.AzureConfiguration) api.Migrator {
	return &Migrator{
		accessor: (&ClusterAccessor{
			configuration: configuration,
			ctx:           context.Background(),
		}).init(),
	}
}
