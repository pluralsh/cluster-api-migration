package worker

import (
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
)

func Defaults() *api.AzureWorker {
	return &api.AzureWorker{
		Replicas: 0,
		Annotations: map[string]string{
			"cluster.x-k8s.io/replicas-managed-by": "external-autoscaler",
		},
		Spec: api.AzureWorkerSpec{
			Mode:              "User",
			SKU:               "Standard_D2s_v3",
			OSDiskSizeGB:      resources.Ptr(int32(50)),
			AdditionalTags:    map[string]*string{},
			AvailabilityZones: []string{"1"},
			NodeLabels: map[string]*string{
				"plural.sh/scalingGroup":    resources.Ptr("medium-sustained-on-demand"),
				"plural.sh/capacityType":    resources.Ptr("ON_DEMAND"),
				"plural.sh/performanceType": resources.Ptr("SUSTAINED"),
			},
			Scaling: &api.ManagedMachinePoolScaling{
				MinSize: 1,
				MaxSize: 5,
			},
			MaxPods:              resources.Ptr(int32(110)),
			OsDiskType:           resources.Ptr("Managed"),
			OSType:               resources.Ptr(api.LinuxOS),
			EnableNodePublicIP:   resources.Ptr(false),
			NodePublicIPPrefixID: resources.Ptr(""),
		},
	}
}
