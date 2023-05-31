package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/migrator"
	"sigs.k8s.io/yaml"
)

const (
	provider = api.ClusterProviderAzure
)

func newConfiguration(provider api.ClusterProvider) *api.Configuration {
	switch provider {
	case api.ClusterProviderGoogle:
		project, region, name := "pluralsh-test-384515", "europe-central2", "gcp-capi"
		credentials, _ := base64.StdEncoding.DecodeString(os.Getenv(api.GCPEncodedCredentialsEnvVar))

		return &api.Configuration{
			GCPConfiguration: &api.GCPConfiguration{
				Credentials: string(credentials),
				Project:     project,
				Region:      region,
				Name:        name,
			},
		}
	case api.ClusterProviderAzure:
		config := api.Configuration{
			AzureConfiguration: &api.AzureConfiguration{
				SubscriptionID: os.Getenv(api.AzureSubscriptionIdEnvVar),
				ResourceGroup:  "plural",
				Name:           "plrltest2",
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
	m := migrator.NewMigrator(provider, newConfiguration(provider))
	res, err := yaml.Marshal(m.Convert())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
}
