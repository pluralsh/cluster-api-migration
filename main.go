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
		return nil
	}

	return nil
}

//func main() {
//	if provider == "aws" {
//		os.Setenv("AWS_ACCESS_KEY_ID", "")
//		os.Setenv("AWS_SECRET_ACCESS_KEY", "")
//		os.Setenv("AWS_SESSION_TOKEN", "")
//		os.Setenv("AWS_REGION", "eu-central-1")
//
//		//cluster, err := aws.GetCluster(context.Background(), "test-aws", "eu-central-1")
//		//if err != nil {
//		//	fmt.Println(err)
//		//}
//
//		//fmt.Printf("cluster %v", cluster)
//	}
//}

func main() {
	m := migrator.NewMigrator(provider, newConfiguration(provider))

	values := m.Convert()
	resources.NewPrinter(values).PrettyPrint()
}
