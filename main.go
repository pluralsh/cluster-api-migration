package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pluralsh/cluster-api-migration/pkg/aws"
	"github.com/pluralsh/cluster-api-migration/pkg/azure"
)

const (
	provider = "aws"
)

func main() {
	if provider == "aws" {
		os.Setenv("AWS_ACCESS_KEY_ID", "")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "")
		os.Setenv("AWS_SESSION_TOKEN", "")
		os.Setenv("AWS_REGION", "eu-central-1")
	
		cluster, err := aws.GetCluster(context.Background(), "test-aws", "eu-central-1")
		if err != nil {
			fmt.Println(err)
		}
	
		fmt.Printf("cluster %v", cluster)
	} else if provider == "azure" {
		// Azure client requires below variables being set:
		// - AZURE_SUBSCRIPTION_ID
		// - AZURE_CLIENT_ID
		// - AZURE_CLIENT_SECRET
		// - AZURE_TENANT_ID

		cluster, err := azure.GetCluster(context.Background(), "plural", "eastus")
		if err != nil {
			fmt.Println(err)
		}
	
		fmt.Printf("cluster %v", cluster)
	}
}
