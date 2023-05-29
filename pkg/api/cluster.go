package api

type Cluster struct {
	Name              string   `json:"name"`
	CIDRBlocks        []string `json:"cidrBlocks"`
	KubernetesVersion string   `json:"kubernetesVersion"`
	CloudSpec         `json:",inline"`
}

type CloudSpec struct {
	AWSCloudSpec *AWSCloudSpec `json:"aws,omitempty"`
}

type Workers struct {
	Defaults    DefaultsWorker `json:"defaults"`
	WorkersSpec `json:",inline"`
}

type WorkersSpec struct {
	AWSWorkers *AWSWorkers `json:"aws,omitempty"`
}

type DefaultsWorker struct {
	AWSDefaultWorker *AWSDefaultWorker `json:"aws,omitempty"`
}
