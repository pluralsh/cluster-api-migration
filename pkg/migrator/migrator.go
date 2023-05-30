package migrator

import (
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp"
)

func NewMigrator(provider api.ClusterProvider, config *api.Configuration) api.Migrator {
	var migrator api.Migrator

	switch provider {
	case api.ClusterProviderGoogle:
		migrator = gcp.NewGCPMigrator(config.GCPConfiguration)
	default:
		return nil
	}

	return migrator
}
