package api

import (
	"fmt"
)

const (
	AzureSubscriptionIdEnvVar   = "AZURE_SUBSCRIPTION_ID"
	AzureClientIdEnvVar         = "AZURE_CLIENT_ID"
	AzureClientSecretEnvVar     = "AZURE_CLIENT_SECRET"
	GCPEncodedCredentialsEnvVar = "GCP_B64ENCODED_CREDENTIALS"
)

type Configuration struct {
	*AWSConfiguration
	*AzureConfiguration
	*GCPConfiguration
}

type AWSConfiguration struct {
}

type AzureConfiguration struct {
	SubscriptionID string
	ResourceGroup  string
	Name           string
}

func (this *AzureConfiguration) Validate() error {
	if len(this.SubscriptionID) == 0 {
		return fmt.Errorf("subscription ID cannot be empty, ensure that %s evironment variable is set", AzureSubscriptionIdEnvVar)
	}

	if len(this.ResourceGroup) == 0 {
		return fmt.Errorf("resource group cannot be empty, ensure that it is set")
	}

	if len(this.Name) == 0 {
		return fmt.Errorf("name cannot be empty, ensure that it is set")
	}

	return nil
}

type GCPConfiguration struct {
	Credentials string
	Project     string
	Region      string
	Name        string
}

type Migrator interface {
	Convert() *Values
}

type ClusterAccessor interface {
	GetCluster() *Cluster
	GetWorkers() *Workers
}
