package migrator

import (
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/aws"
	"github.com/pluralsh/cluster-api-migration/pkg/azure"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp"
)

func NewMigrator(provider api.ClusterProvider, config *api.Configuration) api.Migrator {
	var migrator api.Migrator

	switch provider {
	case api.ClusterProviderGoogle:
		migrator = gcp.NewGCPMigrator(config.GCPConfiguration)
	case api.ClusterProviderAzure:
		migrator = azure.NewAzureMigrator(config.AzureConfiguration)
	case api.ClusterProviderAWS:
		migrator = aws.NewAWSMigrator(config.AWSConfiguration)
	default:
		return nil
	}

	return migrator
}
