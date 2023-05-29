package api

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
}

type Workers struct {
	Defaults    DefaultsWorker `json:"defaults"`
	WorkersSpec `json:",inline"`
}

type WorkersSpec struct {
	AWSWorkers   *AWSWorkers   `json:"aws,omitempty"`
	AzureWorkers *AzureWorkers `json:"azure,omitempty"`
}

type DefaultsWorker struct {
	AWSDefaultWorker   *AWSDefaultWorker   `json:"aws,omitempty"`
	AzureDefaultWorker *AzureDefaultWorker `json:"azure,omitempty"`
}