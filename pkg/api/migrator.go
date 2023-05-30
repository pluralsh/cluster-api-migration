package api

type Configuration struct {
	*AWSConfiguration
	*AzureConfiguration
	*GCPConfiguration
}

type AWSConfiguration struct {
}

type AzureConfiguration struct {
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
	GetCluster() *Values
}
