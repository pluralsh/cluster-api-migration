package api

type GCPReleaseChannel string

const (
	ReleaseChannelStable  = GCPReleaseChannel("stable")
	ReleaseChannelRegular = GCPReleaseChannel("regular")
	ReleaseChannelRapid   = GCPReleaseChannel("rapid")
)

type GCPNetwork struct {
	Name                  string `json:"name"`
	AutoCreateSubnetworks bool   `json:"autoCreateSubnetworks"`
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
	Name                string            `json:"name"`
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

type GCPWorker struct {
	Replicas         int32             `json:"replicas"`
	Scaling          *GCPWorkerScaling `json:"scaling,omitempty"`
	KubernetesLabels *Labels           `json:"kubernetesLabels,omitempty"`
	AdditionalLabels *Labels           `json:"additionalLabels,omitempty"`
	KubernetesTaints *Taints           `json:"kubernetesTaints,omitempty"`
	ProviderIDList   []string          `json:"providerIDList,omitempty"`
}

type GCPWorkerScaling struct {
	MaxCount int32 `json:"maxCount"`
	MinCount int32 `json:"minCount"`
}
