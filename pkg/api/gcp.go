package api

type GCPReleaseChannel string

const (
	ReleaseChannelStable  = GCPReleaseChannel("stable")
	ReleaseChannelRegular = GCPReleaseChannel("regular")
	ReleaseChannelRapid   = GCPReleaseChannel("rapid")
)

// DatapathProvider is the datapath provider selects the implementation of the Kubernetes networking
// model for service resolution and network policy enforcement.
type DatapathProvider string

const (
	// DatapathProvider_UNSPECIFIED is the default value.
	DatapathProvider_UNSPECIFIED DatapathProvider = DatapathProvider("UNSPECIFIED")
	// DatapathProvider_LEGACY_DATAPATH uses the IPTables implementation based on kube-proxy.
	DatapathProvider_LEGACY_DATAPATH DatapathProvider = DatapathProvider("LEGACY_DATAPATH")
	// DatapathProvider_ADVANCED_DATAPATH uses the eBPF based GKE Dataplane V2 with additional features.
	// See the [GKE Dataplane V2 documentation](https://cloud.google.com/kubernetes-engine/docs/how-to/dataplane-v2)
	// for more.
	DatapathProvider_ADVANCED_DATAPATH DatapathProvider = DatapathProvider("ADVANCED_DATAPATH")
)

type GCPNetwork struct {
	Name                  string           `json:"name"`
	AutoCreateSubnetworks bool             `json:"autoCreateSubnetworks"`
	DatapathProvider      DatapathProvider `json:"datapathProvider,omitempty"`
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

type AddonsConfig struct {
	HTTPLoadBalancingEnabled        *bool `json:"httpLoadBalancingEnabled,omitempty"`
	HorizontalPodAutoscalingEnabled *bool `json:"horizontalPodAutoscalingEnabled,omitempty"`
	NetworkPolicyEnabled            *bool `json:"networkPolicyEnabled,omitempty"`
	GcpFilestoreCsiDriverEnabled    *bool `json:"gcpFilestoreCsiDriverEnabled,omitempty"`
}

type GCPCloudSpec struct {
	Project                string             `json:"project"`
	Region                 string             `json:"region"`
	EnableAutopilot        bool               `json:"enableAutopilot"`
	EnableWorkloadIdentity bool               `json:"enableWorkloadIdentity"`
	ReleaseChannel         *GCPReleaseChannel `json:"releaseChannel,omitempty"`
	Network                *GCPNetwork        `json:"network"`
	Subnets                GCPSubnets         `json:"subnets"`
	AdditionalLabels       *Labels            `json:"additionalLabels,omitempty"`
	AddonsConfig           *AddonsConfig      `json:"addonsConfig,omitempty"`
}

type GCPWorkers map[string]*GCPWorker

type GCPWorker struct {
	Replicas          *int32            `json:"replicas,omitempty"`
	KubernetesVersion *string           `json:"kubernetesVersion,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	// IsMultiAZ defines if a node group should be split across the availability zones. If false, will create a node group per AZ
	// +optional
	IsMultiAZ bool          `json:"isMultiAZ,omitempty"`
	Spec      GCPWorkerSpec `json:"spec"`
}

type GCPWorkerSpec struct {
	Scaling          *GCPWorkerScaling    `json:"scaling,omitempty"`
	Management       *GCPWorkerManagement `json:"management,omitempty"`
	KubernetesLabels *Labels              `json:"kubernetesLabels,omitempty"`
	AdditionalLabels *Labels              `json:"additionalLabels,omitempty"`
	KubernetesTaints *Taints              `json:"kubernetesTaints,omitempty"`
	ProviderIDList   []string             `json:"providerIDList,omitempty"`
	MachineType      string               `json:"machineType,omitempty"`
	DiskSizeGb       int32                `json:"diskSizeGb,omitempty"`
	DiskType         string               `json:"diskType,omitempty"`
	ImageType        string               `json:"imageType,omitempty"`
	Preemptible      bool                 `json:"preemptible,omitempty"`
	Spot             bool                 `json:"spot,omitempty"`
}

type GCPWorkerScaling struct {
	MaxCount int32 `json:"maxCount"`
	MinCount int32 `json:"minCount"`
}

type GCPWorkerManagement struct {
	AutoUpgrade bool `json:"autoUpgrade,omitempty"`
	AutoRepair  bool `json:"autoRepair,omitempty"`
}
