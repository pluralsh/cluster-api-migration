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
	Defaults    DefaultsWorkers `json:"defaults"`
	WorkersSpec `json:",inline"`
}

type WorkersSpec struct {
	AWSWorkers *AWSWorkers `json:"aws,omitempty"`
}

type DefaultsWorkers struct {
	AWSDefaultWorkers *AWSDefaultWorkers `json:"aws,omitempty"`
}
