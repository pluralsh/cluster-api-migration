package api

// Following specs should match ones from plural-artifacts repository.

type Tags map[string]string

type IdentityType string

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

type AzureDefaultWorker struct {
	AWSWorker `json:",inline"`
}

type AzureWorkers map[string]AWSWorker

type AzureWorker struct {
	Replicas    int                `json:"replicas"`
	Labels      map[string]*string `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Spec        AWSWorkerSpec      `json:"spec"`
}

type AzureWorkerSpec struct {
	AdditionalTags    Tags                       `json:"additionalTags,omitempty"`
	Mode              string                     `json:"mode"`
	SKU               string                     `json:"sku"`
	OSDiskSizeGB      *int32                     `json:"osDiskSizeGB,omitempty"`
	AvailabilityZones []string                   `json:"availabilityZones,omitempty"`
	NodeLabels        map[string]string          `json:"nodeLabels,omitempty"`
	Taints            Taints                     `json:"taints,omitempty"`
	Scaling           *ManagedMachinePoolScaling `json:"scaling,omitempty"`
	MaxPods           *int32                     `json:"maxPods,omitempty"`
	OsDiskType        *string                    `json:"osDiskType,omitempty"`
	ScaleSetPriority  *string                    `json:"scaleSetPriority,omitempty"`
}
