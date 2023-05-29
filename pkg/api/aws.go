package api

import (
	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// TaintEffect is the effect for a Kubernetes taint.
type TaintEffect string

var (
	// TaintEffectNoSchedule is a taint that indicates that a pod shouldn't be scheduled on a node
	// unless it can tolerate the taint.
	TaintEffectNoSchedule = TaintEffect("no-schedule")
	// TaintEffectNoExecute is a taint that indicates that a pod shouldn't be schedule on a node
	// unless it can tolerate it. And if its already running on the node it will be evicted.
	TaintEffectNoExecute = TaintEffect("no-execute")
	// TaintEffectPreferNoSchedule is a taint that indicates that there is a "preference" that pods shouldn't
	// be scheduled on a node unless it can tolerate the taint. the scheduler will try to avoid placing the pod
	// but it may still run on the node if there is no other option.
	TaintEffectPreferNoSchedule = TaintEffect("prefer-no-schedule")
)

// ManagedMachineAMIType specifies which AWS AMI to use for a managed MachinePool.
type ManagedMachineAMIType string

const (
	// Al2x86_64 is the default AMI type.
	Al2x86_64 ManagedMachineAMIType = "AL2_x86_64"
	// Al2x86_64GPU is the x86-64 GPU AMI type.
	Al2x86_64GPU ManagedMachineAMIType = "AL2_x86_64_GPU"
	// Al2Arm64 is the Arm AMI type.
	Al2Arm64 ManagedMachineAMIType = "AL2_ARM_64"
)

// ManagedMachinePoolCapacityType specifies the capacity type to be used for the managed MachinePool.
type ManagedMachinePoolCapacityType string

const (
	// ManagedMachinePoolCapacityTypeOnDemand is the default capacity type, to launch on-demand instances.
	ManagedMachinePoolCapacityTypeOnDemand ManagedMachinePoolCapacityType = "onDemand"
	// ManagedMachinePoolCapacityTypeSpot is the spot instance capacity type to launch spot instances.
	ManagedMachinePoolCapacityTypeSpot ManagedMachinePoolCapacityType = "spot"
)

// AddonResolution defines the method for resolving parameter conflicts.
type AddonResolution string

var (
	// AddonResolutionOverwrite indicates that if there are parameter conflicts then
	// resolution will be accomplished via overwriting.
	AddonResolutionOverwrite = AddonResolution("overwrite")

	// AddonResolutionNone indicates that if there are parameter conflicts then
	// resolution will not be done and an error will be reported.
	AddonResolutionNone = AddonResolution("none")
)

// Addon represents a EKS addon.
type Addon struct {
	// Name is the name of the addon
	Name string `json:"name"`
	// Version is the version of the addon to use
	Version string `json:"version"`
	// ConflictResolution is used to declare what should happen if there
	// are parameter conflicts. Defaults to none
	ConflictResolution AddonResolution `json:"conflictResolution,omitempty"`
}

// VpcCni specifies configuration related to the VPC CNI.
type VpcCni struct {
	// Disable indicates that the Amazon VPC CNI should be disabled. With EKS clusters the
	// Amazon VPC CNI is automatically installed into the cluster. For clusters where you want
	// to use an alternate CNI this option provides a way to specify that the Amazon VPC CNI
	// should be deleted. You cannot set this to true if you are using the
	// Amazon VPC CNI addon.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`
	// Env defines a list of environment variables to apply to the `aws-node` DaemonSet
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// KubeProxy specifies how the kube-proxy daemonset is managed.
type KubeProxy struct {
	// Disable set to true indicates that kube-proxy should be disabled. With EKS clusters
	// kube-proxy is automatically installed into the cluster. For clusters where you want
	// to use kube-proxy functionality that is provided with an alternate CNI, this option
	// provides a way to specify that the kube-proxy daemonset should be deleted. You cannot
	// set this to true if you are using the Amazon kube-proxy addon.
	Disable bool `json:"disable,omitempty"`
}

// EKSTokenMethod defines the method for obtaining a client token to use when connecting to EKS.
type EKSTokenMethod string

var (
	// EKSTokenMethodIAMAuthenticator indicates that IAM autenticator will be used to get a token.
	EKSTokenMethodIAMAuthenticator = EKSTokenMethod("iam-authenticator")

	// EKSTokenMethodAWSCli indicates that the AWS CLI will be used to get a token
	// Version 1.16.156 or greater is required of the AWS CLI.
	EKSTokenMethodAWSCli = EKSTokenMethod("aws-cli")
)

type AWSCloudSpec struct {
	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName string `json:"sshKeyName,omitempty"`

	// Version defines the desired Kubernetes version. If no version number
	// is supplied then the latest version of Kubernetes that EKS supports
	// will be used.
	// +optional
	Version string `json:"version,omitempty"`

	// RoleName specifies the name of IAM role that gives EKS
	// permission to make API calls. If the role is pre-existing
	// we will treat it as unmanaged and not delete it on
	// deletion. If the EKSEnableIAM feature flag is true
	// and no name is supplied then a role is created.
	// +optional
	RoleName              string                `json:"roleName,omitempty"`
	ControlPlaneEndpoint  clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
	Labels                map[string]string     `json:"labels,omitempty"`
	Addons                []Addon               `json:"addons,omitempty"`
	AssociateOIDCProvider bool                  `json:"associateOIDCProvider,omitempty"`
	Bastion               infrav1.Bastion       `json:"bastion,omitempty"`
	// IdentityRef is a reference to a identity to be used when reconciling the managed control plane.
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec infrav1.NetworkSpec `json:"network,omitempty"`
	// KubeProxy defines managed attributes of the kube-proxy daemonset
	KubeProxy KubeProxy `json:"kubeProxy,omitempty"`
	// VpcCni is used to set configuration options for the VPC CNI plugin
	VpcCni VpcCni `json:"vpcCni,omitempty"`
	// TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
	// iam-authenticator - obtains a client token using iam-authentictor
	// aws-cli - obtains a client token using the AWS CLI
	// Defaults to iam-authenticator
	TokenMethod EKSTokenMethod `json:"tokenMethod,omitempty"`
}

// Taint defines the specs for a Kubernetes taint.
type Taint struct {
	// Effect specifies the effect for the taint
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=no-schedule;no-execute;prefer-no-schedule
	Effect TaintEffect `json:"effect"`
	// Key is the key of the taint
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// Value is the value of the taint
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// Equals is used to test if 2 taints are equal.
func (t *Taint) Equals(other *Taint) bool {
	if t == nil || other == nil {
		return t == other
	}

	return t.Effect == other.Effect &&
		t.Key == other.Key &&
		t.Value == other.Value
}

// Taints is an array of Taints.
type Taints []Taint

// Contains checks for existence of a matching taint.
func (t *Taints) Contains(taint *Taint) bool {
	for _, t := range *t {
		if t.Equals(taint) {
			return true
		}
	}

	return false
}

// UpdateConfig is the configuration options for updating a nodegroup. Only one of MaxUnavailable
// and MaxUnavailablePercentage should be specified.
type UpdateConfig struct {
	// MaxUnavailable is the maximum number of nodes unavailable at once during a version update.
	// Nodes will be updated in parallel. The maximum number is 100.
	// +optional
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	MaxUnavailable *int `json:"maxUnavailable,omitempty"`

	// MaxUnavailablePercentage is the maximum percentage of nodes unavailable during a version update. This
	// percentage of nodes will be updated in parallel, up to 100 nodes at once.
	// +optional
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	MaxUnavailablePercentage *int `json:"maxUnavailablePercentage,omitempty"`
}

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	MinSize int32 `json:"minSize,omitempty"`
	MaxSize int32 `json:"maxSize,omitempty"`
}

type AWSDefaultWorker struct {
	AWSWorker `json:",inline"`
}

type AWSWorkers map[string]AWSWorker

type AWSWorker struct {
	Replicas int `json:"replicas"`
	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]*string `json:"labels,omitempty"`
	// Annotations specifies labels for the Kubernetes node objects
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        AWSWorkerSpec     `json:"spec"`
}

type AWSWorkerSpec struct {
	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]*string `json:"labels,omitempty"`
	// AMIVersion defines the desired AMI release version. If no version number
	// is supplied then the latest version for the Kubernetes version
	// will be used
	// +optional
	AMIVersion string `json:"amiVersion,omitempty"`
	// AMIType defines the AMI type
	// +kubebuilder:validation:Enum:=AL2_x86_64;AL2_x86_64_GPU;AL2_ARM_64;CUSTOM
	// +kubebuilder:default:=AL2_x86_64
	// +optional
	AMIType ManagedMachineAMIType `json:"amiType,omitempty"`
	// CapacityType specifies the capacity type for the ASG behind this pool
	// +kubebuilder:validation:Enum:=onDemand;spot
	// +kubebuilder:default:=onDemand
	// +optional
	CapacityType ManagedMachinePoolCapacityType `json:"capacityType,omitempty"`
	// DiskSize specifies the root disk size
	// +optional
	DiskSize int32 `json:"diskSize,omitempty"`
	// InstanceType specifies the AWS instance type
	// +optional
	InstanceType *string `json:"instanceType,omitempty"`
	// Scaling specifies scaling for the ASG behind this pool
	// +optional
	Scaling *ManagedMachinePoolScaling `json:"scaling,omitempty"`
	// AvailabilityZones is an array of availability zones instances can run in
	AvailabilityZones []string `json:"availabilityZones,omitempty"`
	// SubnetIDs specifies which subnets are used for the
	// auto scaling group of this nodegroup
	// +optional
	SubnetIDs []*string `json:"subnetIDs,omitempty"`
	// Taints specifies the taints to apply to the nodes of the machine pool
	// +optional
	Taints Taints `json:"taints,omitempty"`
	// UpdateConfig holds the optional config to control the behaviour of the update
	// to the nodegroup.
	// +optional
	UpdateConfig *UpdateConfig `json:"updateConfig,omitempty"`
	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`
	// RoleAdditionalPolicies allows you to attach additional polices to
	// the node group role. You must enable the EKSAllowAddRoles
	// feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies []string `json:"roleAdditionalPolicies,omitempty"`
}
