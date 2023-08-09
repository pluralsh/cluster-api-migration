package migrator

import (
	"fmt"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/aws"
	"github.com/pluralsh/cluster-api-migration/pkg/azure"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp"
	"github.com/pluralsh/cluster-api-migration/pkg/kind"
)

func NewMigrator(provider api.ClusterProvider, config *api.Configuration) (api.Migrator, error) {
	switch provider {
	case api.ClusterProviderGoogle:
		return gcp.NewGCPMigrator(config.GCPConfiguration)
	case api.ClusterProviderAzure:
		return azure.NewAzureMigrator(config.AzureConfiguration)
	case api.ClusterProviderAWS:
		return aws.NewAWSMigrator(config.AWSConfiguration)
	case api.ClusterProviderKind:
		return kind.NewKindMigrator(config.KindConfiguration)
	default:
		return nil, fmt.Errorf("unsupported provider")
	}

}
