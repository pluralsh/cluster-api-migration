package worker

import "github.com/pluralsh/cluster-api-migration/pkg/api"

func (this *Workers) defaults() *api.GCPWorker {
	return &api.GCPWorker{
		IsMultiAZ: false,
		Spec: api.GCPWorkerSpec{
			Scaling: &api.GCPWorkerScaling{
				MaxCount: 6,
				MinCount: 3,
			},
			KubernetesLabels: &api.Labels{},
			AdditionalLabels: &api.Labels{},
			KubernetesTaints: &api.Taints{},
			ProviderIDList:   []string{},
			MachineType:      "e2-standard-2",
			DiskSizeGb:       50,
			DiskType:         "pd-standard",
		},
	}
}
