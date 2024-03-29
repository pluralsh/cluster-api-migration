package api

type ClusterProvider string
type ClusterType string

const (
	ClusterProviderAWS   = ClusterProvider("aws")
	ClusterProviderAzure = ClusterProvider("azure")
	ClusterProviderGCP   = ClusterProvider("gcp")
	ClusterProviderKind  = ClusterProvider("kind")

	ClusterTypeManaged = ClusterType("managed")
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
	PodCIDRBlocks     []string `json:"podCidrBlocks,omitempty"`
	ServiceCIDRBlocks []string `json:"serviceCidrBlocks,omitempty"`
	KubernetesVersion string   `json:"kubernetesVersion"`
	CloudSpec         `json:",inline"`
}

type CloudSpec struct {
	AWSCloudSpec   *AWSCloudSpec   `json:"aws,omitempty"`
	AzureCloudSpec *AzureCloudSpec `json:"azure,omitempty"`
	GCPCloudSpec   *GCPCloudSpec   `json:"gcp,omitempty"`
}

type Workers struct {
	WorkersSpec `json:",inline"`
}

type WorkersSpec struct {
	AWSWorkers   *AWSWorkers   `json:"aws,omitempty"`
	AzureWorkers *AzureWorkers `json:"azure,omitempty"`
	GCPWorkers   *GCPWorkers   `json:"gcp,omitempty"`
}
