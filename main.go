package main

import (
	"encoding/base64"
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
		project, region, name := "pluralsh-test-384515", "europe-central2", "gcp-capi"
		credentials, _ := base64.StdEncoding.DecodeString(os.Getenv("GCP_B64ENCODED_CREDENTIALS"))

		return &api.Configuration{
			GCPConfiguration: &api.GCPConfiguration{
				Credentials: string(credentials),
				Project:     project,
				Region:      region,
				Name:        name,
			},
		}
	case api.ClusterProviderAWS:
	case api.ClusterProviderAzure:
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
//	} else if provider == "azure" {
//		// Azure client requires below variables being set:
//		// - AZURE_SUBSCRIPTION_ID
//		// - AZURE_CLIENT_ID
//		// - AZURE_CLIENT_SECRET
//
//		cluster, err := azure.GetCluster(context.Background(), "plrltest2", "plural")
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		fmt.Printf("cluster %v", cluster)
//	}
//}

func main() {
	m := migrator.NewMigrator(provider, newConfiguration(provider))

	values := m.Convert()
	resources.NewPrinter(values).PrettyPrint()
}
