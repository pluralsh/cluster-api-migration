package api

type GCPReleaseChannel string

const (
	ReleaseChannelStable  = GCPReleaseChannel("stable")
	ReleaseChannelRegular = GCPReleaseChannel("regular")
	ReleaseChannelRapid   = GCPReleaseChannel("rapid")
)

type GCPNetwork struct {
	AutoCreateSubnetworks bool `json:"autoCreateSubnetworks"`
}

type GCPSubnetPurpose string

const (
	PurposeInternalHttpsLoadBalancer = GCPSubnetPurpose("INTERNAL_HTTPS_LOAD_BALANCER")
	PurposePrivate                   = GCPSubnetPurpose("PRIVATE")
	PurposePrivateRFC1918            = GCPSubnetPurpose("PRIVATE_RFC_1918")
	PurposePrivateServiceConnect     = GCPSubnetPurpose("PRIVATE_SERVICE_CONNECT")
	PurposeRegionalManagedProxy      = GCPSubnetPurpose("REGIONAL_MANAGED_PROXY")
)

type GCPSubnets []GCPSubnet

type GCPSubnet struct {
	NameSuffix          string            `json:"nameSuffix"`
	CidrBlock           string            `json:"cidrBlock"`
	Description         string            `json:"description"`
	SecondaryCidrBlocks map[string]string `json:"secondaryCidrBlocks"`
	PrivateGoogleAccess bool              `json:"privateGoogleAccess"`
	EnableFlowLogs      bool              `json:"enableFlowLogs"`
	Purpose             GCPSubnetPurpose  `json:"purpose"`
}

type GCPCloudSpec struct {
	Project                string             `json:"project"`
	Region                 string             `json:"region"`
	EnableAutopilot        bool               `json:"enableAutopilot"`
	EnableWorkloadIdentity bool               `json:"enableWorkloadIdentity"`
	ReleaseChannel         *GCPReleaseChannel `json:"releaseChannel,omitempty"`
	Network                *GCPNetwork        `json:"network"`
	Subnets                GCPSubnets         `json:"subnets"`
}

type GCPWorkers map[string]GCPWorker

type GCPDefaultWorker struct {
	GCPWorker `json:",inline"`
}

type GCPWorker struct {
	Replicas         int               `json:"replicas"`
	Scaling          *GCPWorkerScaling `json:"scaling,omitempty"`
	KubernetesLabels *Labels           `json:"kubernetesLabels,omitempty"`
	AdditionalLabels *Labels           `json:"additionalLabels,omitempty"`
	KubernetesTains  *Taints           `json:"kubernetesTains,omitempty"`
	ProviderIDList   []string          `json:"providerIDList,omitempty"`
}

type GCPWorkerScaling struct {
	MaxCount int `json:"maxCount"`
	MinCount int `json:"minCount"`
}