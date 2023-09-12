package api

import (
	"fmt"
)

type Configuration struct {
	*AWSConfiguration
	*AzureConfiguration
	*GCPConfiguration
	*KindConfiguration
}

type AWSConfiguration struct {
	ClusterName string
	Region      string
}

type KindConfiguration struct {
	ClusterName string
}

type AzureConfiguration struct {
	SubscriptionID string
	ResourceGroup  string
	Name           string
}

func (config *AzureConfiguration) Validate() error {
	if len(config.SubscriptionID) == 0 {
		return fmt.Errorf("subscription ID cannot be empty, ensure that it is set")
	}

	if len(config.ResourceGroup) == 0 {
		return fmt.Errorf("resource group cannot be empty, ensure that it is set")
	}

	if len(config.Name) == 0 {
		return fmt.Errorf("name cannot be empty, ensure that it is set")
	}

	return nil
}

type GCPConfiguration struct {
	Project        string
	Region         string
	Name           string
	KubeconfigPath string
}

type Migrator interface {
	Convert() (*Values, error)
	AddTags(tags map[string]string) error
}

type ClusterAccessor interface {
	GetCluster() (*Cluster, error)
	GetWorkers() (*Workers, error)
	AddClusterTags(Tags map[string]string) error
	AddMachinePollsTags(Tags map[string]string) error
	AddVirtualNetworkTags(tags map[string]string) error
}
