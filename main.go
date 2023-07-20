package main

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/migrator"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
)

const (
	provider = api.ClusterProviderAzure
)

func newConfiguration(provider api.ClusterProvider) *api.Configuration {
	switch provider {
	case api.ClusterProviderGoogle:
		kubeconfigPath := os.Getenv("KUBECONFIG")
		credentials, err := base64.StdEncoding.DecodeString(os.Getenv("GCP_B64ENCODED_CREDENTIALS"))
		if err != nil {
			panic(err)
		}

		return &api.Configuration{
			GCPConfiguration: &api.GCPConfiguration{
				Credentials:    string(credentials),
				Project:        "pluralsh-test-384515",
				Region:         "europe-central2",
				Name:           "gcp-capi",
				KubeconfigPath: kubeconfigPath,
			},
		}
	case api.ClusterProviderAzure:
		config := api.Configuration{
			AzureConfiguration: &api.AzureConfiguration{
				SubscriptionID: os.Getenv("AZURE_SUBSCRIPTION_ID"),
				ResourceGroup:  "marcin-aks",
				Name:           "amrc",
			},
		}

		if err := config.Validate(); err != nil {
			log.Fatalln(err)
		}

		return &config
	case api.ClusterProviderAWS:
		config := &api.Configuration{
			AWSConfiguration: &api.AWSConfiguration{
				ClusterName: "lukasz-aws",
				Region:      "eu-central-1",
			},
		}
		return config
	}

	return nil
}

func main() {
	m, err := migrator.NewMigrator(provider, newConfiguration(provider))
	if err != nil {
		log.Fatal(err)
	}

	values, err := m.Convert()
	if err != nil {
		log.Fatal(err)
	}

	resources.NewYAMLPrinter(values).PrettyPrint()
}
