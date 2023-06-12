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
	provider = api.ClusterProviderGoogle
)

func newConfiguration(provider api.ClusterProvider) *api.Configuration {
	switch provider {
	case api.ClusterProviderGoogle:
		credentials, _ := base64.StdEncoding.DecodeString(os.Getenv("GCP_B64ENCODED_CREDENTIALS"))

		return &api.Configuration{
			GCPConfiguration: &api.GCPConfiguration{
				Credentials: string(credentials),
				Project:     "pluralsh-test-384515",
				Region:      "europe-central2",
				Name:        "gcp-capi",
			},
		}
	case api.ClusterProviderAzure:
		config := api.Configuration{
			AzureConfiguration: &api.AzureConfiguration{
				SubscriptionID: os.Getenv("AZURE_SUBSCRIPTION_ID"),
				ResourceGroup:  "aaarg",
				Name:           "aaa",
				// It can be retrieved by using terraform show command in installation repo's bootstap/terraform directory.
				// The path is module.azure-bootstrap.module.aks.tls_private_key.ssh.public_key_openssh.
				SSHPublicKey: os.Getenv("AZURE_B64ENCODED_SSH_PUBLIC_KEY"),
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
