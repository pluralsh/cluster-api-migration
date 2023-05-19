package aws

import (
	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
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

type Cluster struct {
	ControlPlane ControlPlane `json:"controlPlane"`
}

type ControlPlane struct {
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

// EKSTokenMethod defines the method for obtaining a client token to use when connecting to EKS.
type EKSTokenMethod string

var (
	// EKSTokenMethodIAMAuthenticator indicates that IAM autenticator will be used to get a token.
	EKSTokenMethodIAMAuthenticator = EKSTokenMethod("iam-authenticator")

	// EKSTokenMethodAWSCli indicates that the AWS CLI will be used to get a token
	// Version 1.16.156 or greater is required of the AWS CLI.
	EKSTokenMethodAWSCli = EKSTokenMethod("aws-cli")
)

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
