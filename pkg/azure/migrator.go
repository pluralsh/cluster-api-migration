package azure

import (
	"context"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Migrator struct {
	accessor api.ClusterAccessor
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
		Provider: api.ClusterProviderAzure,
		Type:     api.ClusterTypeManaged,
		Cluster:  *c,
		Workers:  *w,
	}, nil
}

func NewAzureMigrator(configuration *api.AzureConfiguration) (api.Migrator, error) {
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
