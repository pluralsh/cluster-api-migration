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

// ControlPlaneLoggingSpec defines what EKS control plane logs that should be enabled.
type ControlPlaneLoggingSpec struct {
	// APIServer indicates if the Kubernetes API Server log (kube-apiserver) shoulkd be enabled
	// +kubebuilder:default=false
	APIServer bool `json:"apiServer"`
	// Audit indicates if the Kubernetes API audit log should be enabled
	Audit bool `json:"audit"`
	// Authenticator indicates if the iam authenticator log should be enabled
	Authenticator bool `json:"authenticator"`
	// ControllerManager indicates if the controller manager (kube-controller-manager) log should be enabled
	ControllerManager bool `json:"controllerManager"`
	// Scheduler indicates if the Kubernetes scheduler (kube-scheduler) log should be enabled
	Scheduler bool `json:"scheduler"`
}

// EncryptionConfig specifies the encryption configuration for the EKS clsuter.
type EncryptionConfig struct {
	// Provider specifies the ARN or alias of the CMK (in AWS KMS)
	Provider string `json:"provider,omitempty"`
	// Resources specifies the resources to be encrypted
	Resources []string `json:"resources,omitempty"`
}

// RoleMapping represents a mapping from a IAM role to Kubernetes users and groups.
type RoleMapping struct {
	// RoleARN is the AWS ARN for the role to map
	// +kubebuilder:validation:MinLength:=31
	RoleARN string `json:"rolearn"`
	// KubernetesMapping holds the RBAC details for the mapping
	KubernetesMapping `json:",inline"`
}

// UserMapping represents a mapping from an IAM user to Kubernetes users and groups.
type UserMapping struct {
	// UserARN is the AWS ARN for the user to map
	// +kubebuilder:validation:MinLength:=31
	UserARN string `json:"userarn"`
	// KubernetesMapping holds the RBAC details for the mapping
	KubernetesMapping `json:",inline"`
}

// KubernetesMapping represents the kubernetes RBAC mapping.
type KubernetesMapping struct {
	// UserName is a kubernetes RBAC user subject
	UserName string `json:"username"`
	// Groups is a list of kubernetes RBAC groups
	Groups []string `json:"groups"`
}

// IAMAuthenticatorConfig represents an aws-iam-authenticator configuration.
type IAMAuthenticatorConfig struct {
	// RoleMappings is a list of role mappings
	// +optional
	RoleMappings []RoleMapping `json:"mapRoles,omitempty"`
	// UserMappings is a list of user mappings
	// +optional
	UserMappings []UserMapping `json:"mapUsers,omitempty"`
}

type OIDCIdentityProviderConfig struct {

	// This is also known as audience. The ID for the client application that makes
	// authentication requests to the OpenID identity provider.
	// +kubebuilder:validation:Required
	ClientID string `json:"clientId,omitempty"`

	// The JWT claim that the provider uses to return your groups.
	// +optional
	GroupsClaim *string `json:"groupsClaim,omitempty"`

	// The prefix that is prepended to group claims to prevent clashes with existing
	// names (such as system: groups). For example, the valueoidc: will create group
	// names like oidc:engineering and oidc:infra.
	// +optional
	GroupsPrefix *string `json:"groupsPrefix,omitempty"`

	// The name of the OIDC provider configuration.
	//
	// IdentityProviderConfigName is a required field
	// +kubebuilder:validation:Required
	IdentityProviderConfigName string `json:"identityProviderConfigName,omitempty"`

	// The URL of the OpenID identity provider that allows the API server to discover
	// public signing keys for verifying tokens. The URL must begin with https://
	// and should correspond to the iss claim in the provider's OIDC ID tokens.
	// Per the OIDC standard, path components are allowed but query parameters are
	// not. Typically the URL consists of only a hostname, like https://server.example.org
	// or https://example.com. This URL should point to the level below .well-known/openid-configuration
	// and must be publicly accessible over the internet.
	//
	// +kubebuilder:validation:Required
	IssuerURL string `json:"issuerUrl,omitempty"`

	// The key value pairs that describe required claims in the identity token.
	// If set, each claim is verified to be present in the token with a matching
	// value. For the maximum number of claims that you can require, see Amazon
	// EKS service quotas (https://docs.aws.amazon.com/eks/latest/userguide/service-quotas.html)
	// in the Amazon EKS User Guide.
	// +optional
	RequiredClaims map[string]string `json:"requiredClaims,omitempty"`

	// The JSON Web Token (JWT) claim to use as the username. The default is sub,
	// which is expected to be a unique identifier of the end user. You can choose
	// other claims, such as email or name, depending on the OpenID identity provider.
	// Claims other than email are prefixed with the issuer URL to prevent naming
	// clashes with other plug-ins.
	// +optional
	UsernameClaim *string `json:"usernameClaim,omitempty"`

	// The prefix that is prepended to username claims to prevent clashes with existing
	// names. If you do not provide this field, and username is a value other than
	// email, the prefix defaults to issuerurl#. You can use the value - to disable
	// all prefixing.
	// +optional
	UsernamePrefix *string `json:"usernamePrefix,omitempty"`

	// tags to apply to oidc identity provider association
	// +optional
	Tags infrav1.Tags `json:"tags,omitempty"`
}

// EndpointAccess specifies how control plane endpoints are accessible.
type EndpointAccess struct {
	// Public controls whether control plane endpoints are publicly accessible
	// +optional
	Public bool `json:"public,omitempty"`
	// PublicCIDRs specifies which blocks can access the public endpoint
	// +optional
	PublicCIDRs []string `json:"publicCIDRs,omitempty"`
	// Private points VPC-internal control plane access to the private endpoint
	// +optional
	Private bool `json:"private,omitempty"`
}

type AWSCloudSpec struct {
	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
	// Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.
	SecondaryCidrBlock string `json:"secondaryCidrBlock"`

	// RoleAdditionalPolicies allows you to attach additional polices to
	// the control plane role. You must enable the EKSAllowAddRoles
	// feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies []string `json:"roleAdditionalPolicies"`

	// EncryptionConfig specifies the encryption configuration for the cluster
	// +optional
	EncryptionConfig EncryptionConfig `json:"encryptionConfig"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	AdditionalTags infrav1.Tags `json:"additionalTags"`

	// IAMAuthenticatorConfig allows the specification of any additional user or role mappings
	// for use when generating the aws-iam-authenticator configuration. If this is nil the
	// default configuration is still generated for the cluster.
	// +optional
	IAMAuthenticatorConfig IAMAuthenticatorConfig `json:"iamAuthenticatorConfig,omitempty"`

	// IdentityProviderconfig is used to specify the oidc provider config
	// to be attached with this eks cluster
	// +optional
	OIDCIdentityProviderConfig OIDCIdentityProviderConfig `json:"oidcIdentityProviderConfig,omitempty"`

	// Logging specifies which EKS Cluster logs should be enabled. Entries for
	// each of the enabled logs will be sent to CloudWatch
	// +optional
	Logging ControlPlaneLoggingSpec `json:"logging,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName string `json:"sshKeyName,omitempty"`

	// Version defines the desired Kubernetes version. If no version number
	// is supplied then the latest version of Kubernetes that EKS supports
	// will be used.
	// +optional
	Version string `json:"version,omitempty"`

	// Endpoints specifies access to this cluster's control plane endpoints
	// +optional
	EndpointAccess EndpointAccess `json:"endpointAccess"`

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

type AWSWorkers map[string]*AWSWorker

type AWSWorker struct {
	Replicas int `json:"replicas"`
	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]*string `json:"labels,omitempty"`
	// Annotations specifies labels for the Kubernetes node objects
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// IsMultiAZ defines if a node group should be split across the availability zones. If false, will create a node group per AZ
	// +optional
	IsMultiAZ bool          `json:"isMultiAZ,omitempty"`
	Spec      AWSWorkerSpec `json:"spec"`
}

type AWSWorkerSpec struct {
	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]*string `json:"labels"`
	// AMIVersion defines the desired AMI release version. If no version number
	// is supplied then the latest version for the Kubernetes version
	// will be used
	// +optional
	AMIVersion string `json:"amiVersion"`
	// AMIType defines the AMI type
	// +optional
	AMIType ManagedMachineAMIType `json:"amiType,omitempty"`
	// CapacityType specifies the capacity type for the ASG behind this pool
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
	AdditionalTags infrav1.Tags `json:"additionalTags"`
	// RoleAdditionalPolicies allows you to attach additional polices to
	// the node group role. You must enable the EKSAllowAddRoles
	// feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies []string `json:"roleAdditionalPolicies,omitempty"`
}
