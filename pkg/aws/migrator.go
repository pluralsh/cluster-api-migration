package aws

import (
	"context"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Migrator struct {
	accessor api.ClusterAccessor
}

func (m Migrator) Convert() *api.Values {
	c := m.accessor.GetCluster()
	w := m.accessor.GetWorkers()

	return &api.Values{
		Provider: api.ClusterProviderAzure,
		Type:     api.ClusterTypeManaged,
		Cluster:  *c,
		Workers:  *w,
	}
}

func NewAWSMigrator(configuration *api.AWSConfiguration) api.Migrator {
	return &Migrator{
		accessor: (&ClusterAccessor{
			configuration: configuration,
			ctx:           context.Background(),
		}).init(),
	}
}
