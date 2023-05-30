package gcp

import (
	"context"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
)

type Migrator struct {
	accessor api.ClusterAccessor
}

func (this *Migrator) Convert() *api.Values {
	c := this.accessor.GetCluster()

	resources.NewPrinter(c).PrettyPrint()

	return nil
}

func NewGCPMigrator(configuration *api.GCPConfiguration) api.Migrator {
	return &Migrator{
		accessor: (&ClusterAccessor{
			configuration: configuration,
			ctx:           context.Background(),
		}).init(),
	}
}
