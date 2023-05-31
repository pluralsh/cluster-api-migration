package worker

import (
	"cloud.google.com/go/container/apiv1/containerpb"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Workers struct {
	*containerpb.Cluster
}

func (this *Workers) toGCPWorkers() *api.GCPWorkers {
	workers := api.GCPWorkers{}

	for _, nodePool := range this.Cluster.NodePools {
		workers[nodePool.Name] = this.toGCPWorker(nodePool)
	}

	return &workers
}

func (this *Workers) toGCPWorker(nodePool *containerpb.NodePool) api.GCPWorker {
	var autoscaling *api.GCPWorkerScaling
	if nodePool.GetAutoscaling() != nil {
		autoscaling = &api.GCPWorkerScaling{
			MaxCount: nodePool.Autoscaling.MaxNodeCount,
			MinCount: nodePool.Autoscaling.MinNodeCount,
		}
	}

	return api.GCPWorker{
		Replicas:         nodePool.InitialNodeCount,
		Scaling:          autoscaling,
		KubernetesLabels: nil,
		AdditionalLabels: nil,
		KubernetesTaints: nil,
		ProviderIDList:   nil,
	}
}

func (this *Workers) Convert() *api.Workers {
	return &api.Workers{
		Defaults: api.DefaultsWorker{
			GCPDefaultWorker: this.defaults(),
		},
		WorkersSpec: api.WorkersSpec{
			GCPWorkers: this.toGCPWorkers(),
		},
	}
}

func NewGCPWorkers(cluster *containerpb.Cluster) *Workers {
	return &Workers{
		Cluster: cluster,
	}
}
