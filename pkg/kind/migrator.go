package kind

import (
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Migrator struct {
	accessor api.ClusterAccessor
}

func (m Migrator) AddTags(tags map[string]string) error {
	return nil
}

func (m Migrator) Convert() (*api.Values, error) {
	return &api.Values{
		Provider: api.ClusterProviderKind,
		Type:     api.ClusterTypeManaged,
		Cluster:  api.Cluster{},
		Workers:  api.Workers{},
	}, nil
}

func NewKindMigrator(configuration *api.KindConfiguration) (api.Migrator, error) {
	a, err := (&ClusterAccessor{
		configuration: configuration,
	}).init()
	if err != nil {
		return nil, err
	}
	return &Migrator{
		accessor: a,
	}, nil
}
