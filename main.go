package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pluralsh/cluster-api-migration/pkg/aws"
)

func main() {

	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")
	os.Setenv("AWS_SESSION_TOKEN", "")
	os.Setenv("AWS_REGION", "eu-central-1")

	cluster, err := aws.GetCluster(context.Background(), "test-aws", "eu-central-1")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("cluster %v", cluster)

}
