package api

type ClusterProvider string
type ClusterType string

const (
	ClusterProviderAWS    = ClusterProvider("aws")
	ClusterProviderAzure  = ClusterProvider("azure")
	ClusterProviderGoogle = ClusterProvider("google")

	ClusterTypeManaged   = ClusterType("managed")
	ClusterTypeUnmanaged = ClusterType("unmanaged")
)

type Values struct {
	Provider ClusterProvider `json:"provider"`
	Type     ClusterType     `json:"type"`
	Cluster  Cluster         `json:"cluster"`
	Workers  Workers         `json:"workers"`
}

type ClusterAPI struct {
	Provider string  `json:"provider"`
	Type     string  `json:"type"`
	Cluster  Cluster `json:"cluster"`
	Workers  Workers `json:"workers"`
}

type Cluster struct {
	Name              string   `json:"name"`
	CIDRBlocks        []string `json:"cidrBlocks"`
	KubernetesVersion string   `json:"kubernetesVersion"`
	CloudSpec         `json:",inline"`
}

type CloudSpec struct {
	AWSCloudSpec   *AWSCloudSpec   `json:"aws,omitempty"`
	AzureCloudSpec *AzureCloudSpec `json:"azure,omitempty"`
	GCPCloudSpec   *GCPCloudSpec   `json:"google,omitempty"`
}

type Workers struct {
	Defaults    DefaultsWorker `json:"defaults"`
	WorkersSpec `json:",inline"`
}

type WorkersSpec struct {
	AWSWorkers   *AWSWorkers   `json:"aws,omitempty"`
	AzureWorkers *AzureWorkers `json:"azure,omitempty"`
	GCPWorkers   *GCPWorkers   `json:"google,omitempty"`
}

type DefaultsWorker struct {
	AWSDefaultWorker   *AWSWorker   `json:"aws,omitempty"`
	AzureDefaultWorker *AzureWorker `json:"azure,omitempty"`
	GCPDefaultWorker   *GCPDefaultWorker   `json:"google,omitempty"`
}
