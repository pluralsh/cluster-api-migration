package api

import (
	corev1 "k8s.io/api/core/v1"
)

type Tags map[string]string

type IdentityType string

type AzureCloudSpec struct {
	ClusterIdentityName string                 `json:"clusterIdentityName,omitempty"`
	ClusterIdentityType IdentityType           `json:"clusterIdentityType,omitempty"`
	ClientID            string                 `json:"clientID"`
	ClientSecret        corev1.SecretReference `json:"clientSecret,omitempty"`
	ClientSecretName    string                 `json:"clientSecretName,omitempty"`
	ResourceID          string                 `json:"resourceID,omitempty"`
	TenantID            string                 `json:"tenantID"`
	SubscriptionID      string                 `json:"subscriptionID"`
	Location            string                 `json:"location"`
	ResourceGroupName   string                 `json:"resourceGroupName"`
	SSHPublicKey        string                 `json:"sshPublicKey"`
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
