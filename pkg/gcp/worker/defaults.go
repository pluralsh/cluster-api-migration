package worker

import "github.com/pluralsh/cluster-api-migration/pkg/api"

func (this *Workers) defaults() *api.GCPWorker {
	return &api.GCPWorker{
		Scaling: &api.GCPWorkerScaling{
			MaxCount: 6,
			MinCount: 3,
		},
		KubernetesLabels: &api.Labels{},
		AdditionalLabels: &api.Labels{},
		KubernetesTaints: &api.Taints{},
		ProviderIDList:   []string{},
	}
}
