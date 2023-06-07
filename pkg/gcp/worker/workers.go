package worker

import (
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"

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
		Scaling:          autoscaling,
		KubernetesLabels: this.kubernetesLabels(nodePool),
		AdditionalLabels: this.additionalLabels(nodePool),
		KubernetesTaints: this.kubernetesTaints(nodePool),
		ProviderIDList:   []string{},
	}
}

func (this *Workers) kubernetesLabels(nodePool *containerpb.NodePool) *api.Labels {
	if nodePool == nil || nodePool.Config == nil {
		return nil
	}

	return resources.Ptr(api.Labels(nodePool.Config.Labels))
}

func (this *Workers) additionalLabels(nodePool *containerpb.NodePool) *api.Labels {
	if nodePool == nil || nodePool.Config == nil {
		return nil
	}

	return resources.Ptr(api.Labels(nodePool.Config.Metadata))
}

func (this *Workers) kubernetesTaints(nodePool *containerpb.NodePool) *api.Taints {
	if nodePool == nil || nodePool.Config == nil {
		return nil
	}

	return this.toTaints(nodePool.Config.Taints)
}

func (this *Workers) toTaints(taints []*containerpb.NodeTaint) *api.Taints {
	result := make([]api.Taint, 0)
	for _, taint := range taints {
		result = append(result, api.Taint{
			Effect: this.toTaintEffect(taint.Effect),
			Key:    taint.Key,
			Value:  taint.Value,
		})
	}

	return resources.Ptr(api.Taints(result))
}

func (this *Workers) toTaintEffect(effect containerpb.NodeTaint_Effect) api.TaintEffect {
	switch effect {
	case containerpb.NodeTaint_NO_SCHEDULE:
		return "NoSchedule"
	case containerpb.NodeTaint_NO_EXECUTE:
		return "NoExecute"
	case containerpb.NodeTaint_PREFER_NO_SCHEDULE:
		return "PreferNoSchedule"
	default:
		return ""
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
