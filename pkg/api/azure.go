package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Following specs should match ones from plural-artifacts repository.

type AzureCloudSpec struct {
	// Name of AzureClusterIdentity to be used when reconciling this cluster.
	ClusterIdentityName string `json:"clusterIdentityName"`

	// Type of Azure Identity used.
	// One of: ServicePrincipal, ServicePrincipalCertificate, UserAssignedMSI or ManualServicePrincipal.
	ClusterIdentityType IdentityType `json:"clusterIdentityType"`

	// AllowedNamespaces is used to identify the namespaces the clusters are allowed to use the identity from.
	// Namespaces can be selected either using an array of namespaces or with label selector.
	// An empty allowedNamespaces object indicates that AzureClusters can use this identity from any namespace.
	// If this object is nil, no namespaces will be allowed (default behaviour, if this field is not provided)
	// A namespace should be either in the NamespaceList or match with Selector to use the identity.
	AllowedNamespaces *AllowedNamespaces `json:"allowedNamespaces"`

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

	// NodeResourceGroupName is the name of the resource group
	// containing cluster IaaS resources. Will be populated to default
	// in webhook.
	NodeResourceGroupName string `json:"nodeResourceGroupName,omitempty"`

	// VirtualNetwork describes the vnet for the AKS cluster. Will be created if it does not exist.
	VirtualNetwork ManagedControlPlaneVirtualNetwork `json:"virtualNetwork,omitempty"`

	// NetworkPlugin used for building Kubernetes network.
	// +kubebuilder:validation:Enum=azure;kubenet
	// +optional
	NetworkPlugin *string `json:"networkPlugin,omitempty"`

	// NetworkPolicy used for building Kubernetes network.
	// +kubebuilder:validation:Enum=azure;calico
	// +optional
	NetworkPolicy *string `json:"networkPolicy,omitempty"`

	// Outbound configuration used by Nodes.
	// +kubebuilder:validation:Enum=loadBalancer;managedNATGateway;userAssignedNATGateway;userDefinedRouting
	// +optional
	OutboundType *string `json:"outboundType,omitempty"`

	// DNSServiceIP is an IP address assigned to the Kubernetes DNS service.
	// It must be within the Kubernetes service address range specified in serviceCidr.
	// +optional
	DNSServiceIP *string `json:"dnsServiceIP,omitempty"`

	// SKU is the SKU of the AKS to be provisioned.
	// +optional
	SKU *AKSSku `json:"sku,omitempty"`

	// LoadBalancerSKU is the SKU of the loadBalancer to be provisioned.
	// +kubebuilder:validation:Enum=Basic;Standard
	// +optional
	LoadBalancerSKU *string `json:"loadBalancerSKU,omitempty"`

	// String literal containing an SSH public key base64 encoded.
	SSHPublicKey string `json:"sshPublicKey"`

	// LoadBalancerProfile is the profile of the cluster load balancer.
	// +optional
	LoadBalancerProfile *LoadBalancerProfile `json:"loadBalancerProfile,omitempty"`

	// APIServerAccessProfile is the access profile for AKS API server.
	// +optional
	APIServerAccessProfile *APIServerAccessProfile `json:"apiServerAccessProfile,omitempty"`

	// AutoscalerProfile is the parameters to be applied to the cluster-autoscaler when enabled
	// +optional
	AutoScalerProfile *AutoScalerProfile `json:"autoscalerProfile,omitempty"`

	// AadProfile is Azure Active Directory configuration to integrate with AKS for aad authentication.
	// +optional
	AADProfile *AADProfile `json:"aadProfile,omitempty"`

	// AddonProfiles are the profiles of managed cluster add-on.
	// +optional
	AddonProfiles []AddonProfile `json:"addonProfiles,omitempty"`
}

type AzureWorkers map[string]AzureWorker

type AzureWorker struct {
	Replicas    int               `json:"replicas"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        AzureWorkerSpec   `json:"spec"`
}

type AzureWorkerSpec struct {
	AdditionalTags       map[string]*string         `json:"additionalTags"`
	Mode                 string                     `json:"mode"`
	SKU                  string                     `json:"sku"`
	OSDiskSizeGB         *int32                     `json:"osDiskSizeGB,omitempty"`
	AvailabilityZones    []*string                  `json:"availabilityZones,omitempty"`
	NodeLabels           map[string]*string         `json:"nodeLabels"`
	Taints               Taints                     `json:"taints,omitempty"`
	Scaling              *ManagedMachinePoolScaling `json:"scaling,omitempty"`
	MaxPods              *int32                     `json:"maxPods,omitempty"`
	OsDiskType           *string                    `json:"osDiskType,omitempty"`
	OSType               *string                    `json:"osType,omitempty"`
	EnableNodePublicIP   *bool                      `json:"enableNodePublicIP,omitempty"`
	NodePublicIPPrefixID *string                    `json:"nodePublicIPPrefixID,omitempty"`
	KubeletConfig        *KubeletConfig             `json:"kubeletConfig,omitempty"`
	LinuxOSConfig        *LinuxOSConfig             `json:"linuxOSConfig,omitempty"`
	ScaleSetPriority     *string                    `json:"scaleSetPriority,omitempty"`
}

type AllowedNamespaces struct {
	NamespaceList []string              `json:"list,omitempty"`
	Selector      *metav1.LabelSelector `json:"selector,omitempty"`
}

type ManagedControlPlaneVirtualNetwork struct {
	Name          string                    `json:"name"`
	CIDRBlock     string                    `json:"cidrBlock"`
	Subnet        ManagedControlPlaneSubnet `json:"subnet,omitempty"`
	ResourceGroup string                    `json:"resourceGroup,omitempty"`
}

type ManagedControlPlaneSubnet struct {
	Name             string           `json:"name"`
	CIDRBlock        string           `json:"cidrBlock"`
	ServiceEndpoints ServiceEndpoints `json:"serviceEndpoints,omitempty"`
	PrivateEndpoints PrivateEndpoints `json:"privateEndpoints,omitempty"`
}

type ServiceEndpoints []ServiceEndpointSpec

type ServiceEndpointSpec struct {
	Service   string   `json:"service"`
	Locations []string `json:"locations"`
}

type PrivateEndpoints []PrivateEndpointSpec

type PrivateEndpointSpec struct {
	Name                          string                         `json:"name"`
	Location                      string                         `json:"location,omitempty"`
	PrivateLinkServiceConnections []PrivateLinkServiceConnection `json:"privateLinkServiceConnections,omitempty"`
	CustomNetworkInterfaceName    string                         `json:"customNetworkInterfaceName,omitempty"`
	PrivateIPAddresses            []string                       `json:"privateIPAddresses,omitempty"`
	ApplicationSecurityGroups     []string                       `json:"applicationSecurityGroups,omitempty"`
	ManualApproval                bool                           `json:"manualApproval,omitempty"`
}

type PrivateLinkServiceConnection struct {
	Name                 string   `json:"name,omitempty"`
	PrivateLinkServiceID string   `json:"privateLinkServiceID,omitempty"`
	GroupIDs             []string `json:"groupIDs,omitempty"`
	RequestMessage       string   `json:"requestMessage,omitempty"`
}

type AADProfile struct {
	Managed             bool     `json:"managed"`
	AdminGroupObjectIDs []string `json:"adminGroupObjectIDs"`
}

type AddonProfile struct {
	Name    string             `json:"name"`
	Config  map[string]*string `json:"config,omitempty"`
	Enabled bool               `json:"enabled"`
}

type AutoScalerProfile struct {
	BalanceSimilarNodeGroups      *string `json:"balanceSimilarNodeGroups,omitempty"`
	Expander                      *string `json:"expander,omitempty"`
	MaxEmptyBulkDelete            *string `json:"maxEmptyBulkDelete,omitempty"`
	MaxGracefulTerminationSec     *string `json:"maxGracefulTerminationSec,omitempty"`
	MaxNodeProvisionTime          *string `json:"maxNodeProvisionTime,omitempty"`
	MaxTotalUnreadyPercentage     *string `json:"maxTotalUnreadyPercentage,omitempty"`
	NewPodScaleUpDelay            *string `json:"newPodScaleUpDelay,omitempty"`
	OkTotalUnreadyCount           *string `json:"okTotalUnreadyCount,omitempty"`
	ScanInterval                  *string `json:"scanInterval,omitempty"`
	ScaleDownDelayAfterAdd        *string `json:"scaleDownDelayAfterAdd,omitempty"`
	ScaleDownDelayAfterDelete     *string `json:"scaleDownDelayAfterDelete,omitempty"`
	ScaleDownDelayAfterFailure    *string `json:"scaleDownDelayAfterFailure,omitempty"`
	ScaleDownUnneededTime         *string `json:"scaleDownUnneededTime,omitempty"`
	ScaleDownUnreadyTime          *string `json:"scaleDownUnreadyTime,omitempty"`
	ScaleDownUtilizationThreshold *string `json:"scaleDownUtilizationThreshold,omitempty"`
	SkipNodesWithLocalStorage     *string `json:"skipNodesWithLocalStorage,omitempty"`
	SkipNodesWithSystemPods       *string `json:"skipNodesWithSystemPods,omitempty"`
}

type AKSSku struct {
	Tier *string `json:"tier"`
}

type LoadBalancerProfile struct {
	ManagedOutboundIPs     *int32   `json:"managedOutboundIPs,omitempty"`
	OutboundIPPrefixes     []string `json:"outboundIPPrefixes,omitempty"`
	OutboundIPs            []string `json:"outboundIPs,omitempty"`
	AllocatedOutboundPorts *int32   `json:"allocatedOutboundPorts,omitempty"`
	IdleTimeoutInMinutes   *int32   `json:"idleTimeoutInMinutes,omitempty"`
}

type APIServerAccessProfile struct {
	AuthorizedIPRanges             []*string `json:"authorizedIPRanges,omitempty"`
	EnablePrivateCluster           *bool     `json:"enablePrivateCluster,omitempty"`
	PrivateDNSZone                 *string   `json:"privateDNSZone,omitempty"`
	EnablePrivateClusterPublicFQDN *bool     `json:"enablePrivateClusterPublicFQDN,omitempty"`
}

type IdentityType string

type KubeletConfig struct {
	CPUManagerPolicy      *CPUManagerPolicy      `json:"cpuManagerPolicy,omitempty"`
	CPUCfsQuota           *bool                  `json:"cpuCfsQuota,omitempty"`
	CPUCfsQuotaPeriod     *string                `json:"cpuCfsQuotaPeriod,omitempty"`
	ImageGcHighThreshold  *int32                 `json:"imageGcHighThreshold,omitempty"`
	ImageGcLowThreshold   *int32                 `json:"imageGcLowThreshold,omitempty"`
	TopologyManagerPolicy *TopologyManagerPolicy `json:"topologyManagerPolicy,omitempty"`
	AllowedUnsafeSysctls  []string               `json:"allowedUnsafeSysctls,omitempty"`
	FailSwapOn            *bool                  `json:"failSwapOn,omitempty"`
	ContainerLogMaxSizeMB *int32                 `json:"containerLogMaxSizeMB,omitempty"`
	ContainerLogMaxFiles  *int32                 `json:"containerLogMaxFiles,omitempty"`
	PodMaxPids            *int32                 `json:"podMaxPids,omitempty"`
}

type CPUManagerPolicy string

const (
	CPUManagerPolicyNone   CPUManagerPolicy = "none"
	CPUManagerPolicyStatic CPUManagerPolicy = "static"
)

type TopologyManagerPolicy string

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
