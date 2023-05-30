package api

// Following specs should match ones from plural-artifacts repository.

type AzureCloudSpec struct {
	// Name of AzureClusterIdentity to be used when reconciling this cluster.
	ClusterIdentityName string `json:"clusterIdentityName,omitempty"`

	// Type of Azure Identity used.
	// One of: ServicePrincipal, ServicePrincipalCertificate, UserAssignedMSI or ManualServicePrincipal.
	ClusterIdentityType IdentityType `json:"clusterIdentityType,omitempty"`

	// Service Principal client ID. Both Service Principal and User Assigned MSI can use this field.
	ClientID string `json:"clientID"`

	// Service Principal password.
	ClientSecret string `json:"clientSecret,omitempty"`

	// Name of Secret containing clientSecret.
	ClientSecretName string `json:"clientSecretName,omitempty"`

	// Azure resource ID for the User Assigned MSI resource.
	ResourceID string `json:"resourceID,omitempty"`

	// Service Principal primary tenant ID.
	TenantID string `json:"tenantID"`

	// GUID of the Azure subscription to hold this cluster.
	SubscriptionID string `json:"subscriptionID"`

	// String matching one of the canonical Azure region names.
	// Examples: "westus2", "eastus".
	Location string `json:"location"`

	// Name of the Azure resource group for this AKS Cluster.
	ResourceGroupName string `json:"resourceGroupName"`

	// String literal containing an SSH public key base64 encoded.
	SSHPublicKey string `json:"sshPublicKey"`
}

type AzureWorkers map[string]AzureWorker

type AzureWorker struct {
	Replicas    int                `json:"replicas"`
	Labels      map[string]*string `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Spec        AzureWorkerSpec    `json:"spec"`
}

type AzureWorkerSpec struct {
	AdditionalTags       Tags                       `json:"additionalTags,omitempty"`
	Mode                 string                     `json:"mode"`
	SKU                  string                     `json:"sku"`
	OSDiskSizeGB         *int32                     `json:"osDiskSizeGB,omitempty"`
	AvailabilityZones    []string                   `json:"availabilityZones,omitempty"`
	NodeLabels           map[string]string          `json:"nodeLabels,omitempty"`
	Taints               Taints                     `json:"taints,omitempty"`
	Scaling              *ManagedMachinePoolScaling `json:"scaling,omitempty"`
	MaxPods              *int32                     `json:"maxPods,omitempty"`
	OsDiskType           *string                    `json:"osDiskType,omitempty"`
	OSType               *string                    `json:"osType,omitempty"`
	EnableNodePublicIP   *bool                      `json:"enableNodePublicIP,omitempty"`
	NodePublicIPPrefixID *string                    `json:"nodePublicIPPrefixID,omitempty"`
	KubeletDiskType      *KubeletDiskType           `json:"kubeletDiskType,omitempty"`
	LinuxOSConfig        *LinuxOSConfig             `json:"linuxOSConfig,omitempty"`
	ScaleSetPriority     *string                    `json:"scaleSetPriority,omitempty"`
}

type Tags map[string]string

type IdentityType string

type KubeletDiskType string

const (
	KubeletDiskTypeOS        KubeletDiskType = "OS"
	KubeletDiskTypeTemporary KubeletDiskType = "Temporary"
)

type LinuxOSConfig struct {
	SwapFileSizeMB             *int32                     `json:"swapFileSizeMB,omitempty"`
	Sysctls                    *SysctlConfig              `json:"sysctls,omitempty"`
	TransparentHugePageDefrag  *TransparentHugePageOption `json:"transparentHugePageDefrag,omitempty"`
	TransparentHugePageEnabled *TransparentHugePageOption `json:"transparentHugePageEnabled,omitempty"`
}

type SysctlConfig struct {
	FsAioMaxNr                     *int32  `json:"fsAioMaxNr,omitempty"`
	FsFileMax                      *int32  `json:"fsFileMax,omitempty"`
	FsInotifyMaxUserWatches        *int32  `json:"fsInotifyMaxUserWatches,omitempty"`
	FsNrOpen                       *int32  `json:"fsNrOpen,omitempty"`
	KernelThreadsMax               *int32  `json:"kernelThreadsMax,omitempty"`
	NetCoreNetdevMaxBacklog        *int32  `json:"netCoreNetdevMaxBacklog,omitempty"`
	NetCoreOptmemMax               *int32  `json:"netCoreOptmemMax,omitempty"`
	NetCoreRmemDefault             *int32  `json:"netCoreRmemDefault,omitempty"`
	NetCoreRmemMax                 *int32  `json:"netCoreRmemMax,omitempty"`
	NetCoreSomaxconn               *int32  `json:"netCoreSomaxconn,omitempty"`
	NetCoreWmemDefault             *int32  `json:"netCoreWmemDefault,omitempty"`
	NetCoreWmemMax                 *int32  `json:"netCoreWmemMax,omitempty"`
	NetIpv4IPLocalPortRange        *string `json:"netIpv4IPLocalPortRange,omitempty"`
	NetIpv4NeighDefaultGcThresh1   *int32  `json:"netIpv4NeighDefaultGcThresh1,omitempty"`
	NetIpv4NeighDefaultGcThresh2   *int32  `json:"netIpv4NeighDefaultGcThresh2,omitempty"`
	NetIpv4NeighDefaultGcThresh3   *int32  `json:"netIpv4NeighDefaultGcThresh3,omitempty"`
	NetIpv4TCPFinTimeout           *int32  `json:"netIpv4TCPFinTimeout,omitempty"`
	NetIpv4TCPKeepaliveProbes      *int32  `json:"netIpv4TCPKeepaliveProbes,omitempty"`
	NetIpv4TCPKeepaliveTime        *int32  `json:"netIpv4TCPKeepaliveTime,omitempty"`
	NetIpv4TCPMaxSynBacklog        *int32  `json:"netIpv4TCPMaxSynBacklog,omitempty"`
	NetIpv4TCPMaxTwBuckets         *int32  `json:"netIpv4TCPMaxTwBuckets,omitempty"`
	NetIpv4TCPTwReuse              *bool   `json:"netIpv4TCPTwReuse,omitempty"`
	NetIpv4TCPkeepaliveIntvl       *int32  `json:"netIpv4TCPkeepaliveIntvl,omitempty"`
	NetNetfilterNfConntrackBuckets *int32  `json:"netNetfilterNfConntrackBuckets,omitempty"`
	NetNetfilterNfConntrackMax     *int32  `json:"netNetfilterNfConntrackMax,omitempty"`
	VMMaxMapCount                  *int32  `json:"vmMaxMapCount,omitempty"`
	VMSwappiness                   *int32  `json:"vmSwappiness,omitempty"`
	VMVfsCachePressure             *int32  `json:"vmVfsCachePressure,omitempty"`
}

type TransparentHugePageOption string

const (
	TransparentHugePageOptionAlways       TransparentHugePageOption = "always"
	TransparentHugePageOptionMadvise      TransparentHugePageOption = "madvise"
	TransparentHugePageOptionNever        TransparentHugePageOption = "never"
	TransparentHugePageOptionDefer        TransparentHugePageOption = "defer"
	TransparentHugePageOptionDeferMadvise TransparentHugePageOption = "defer+madvise"
)
