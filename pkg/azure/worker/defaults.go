package worker

import "github.com/pluralsh/cluster-api-migration/pkg/api"

func AzureWorkerDefaults() *api.AzureWorker {
	osDiskSizeGB := int32(50)
	osDiskType := "Managed"
	osType := api.LinuxOS
	maxPods := int32(110)
	enableNodePublicIP := false
	nodePublicIPPrefixID := ""
	scaleSetPriority := "Regular"

	return &api.AzureWorker{
		Replicas: 0,
		Annotations: map[string]string{
			"cluster.x-k8s.io/replicas-managed-by": "external-autoscaler",
		},
		Spec: api.AzureWorkerSpec{
			Mode:              api.NodePoolModeUser,
			SKU:               "Standard_D2s_v3",
			OSDiskSizeGB:      &osDiskSizeGB,
			AvailabilityZones: []string{"1"},
			NodeLabels: map[string]string{
				"plural.sh/scalingGroup":    "medium-sustained-on-demand",
				"plural.sh/capacityType":    "ON_DEMAND",
				"plural.sh/performanceType": "SUSTAINED",
			},
			Scaling: &api.ManagedMachinePoolScaling{
				MinSize: 1,
				MaxSize: 5,
			},
			MaxPods:              &maxPods,
			OsDiskType:           &osDiskType,
			OSType:               &osType,
			EnableNodePublicIP:   &enableNodePublicIP,
			NodePublicIPPrefixID: &nodePublicIPPrefixID,
			ScaleSetPriority:     &scaleSetPriority,
		},
	}
}
