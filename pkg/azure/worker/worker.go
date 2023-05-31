package worker

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v2"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

type Workers struct {
	Cluster        *armcontainerservice.ManagedCluster
	ResourceGroup  string
	SubscriptionID string
}

func (workers *Workers) Workers() *api.AzureWorkers {
	result := api.AzureWorkers{}

	for _, agentPool := range workers.Cluster.Properties.AgentPoolProfiles {
		result[*agentPool.Name] = workers.Worker(agentPool)
	}

	return &result
}

func (workers *Workers) Worker(agentPool *armcontainerservice.ManagedClusterAgentPoolProfile) api.AzureWorker {
	worker := api.AzureWorker{
		Replicas:    int(*agentPool.Count),
		Annotations: map[string]string{}, // TODO:
		Spec: api.AzureWorkerSpec{
			AdditionalTags:    nil, // TODO:
			Mode:              string(*agentPool.Mode),
			SKU:               *agentPool.VMSize,
			OSDiskSizeGB:      agentPool.OSDiskSizeGB,
			AvailabilityZones: agentPool.AvailabilityZones,
			NodeLabels:        agentPool.NodeLabels,
			Scaling: &api.ManagedMachinePoolScaling{
				MinSize: *agentPool.MinCount,
				MaxSize: *agentPool.MaxCount,
			},
			MaxPods:              agentPool.MaxPods,
			OsDiskType:           (*string)(agentPool.OSDiskType),
			OSType:               (*string)(agentPool.OSType),
			EnableNodePublicIP:   agentPool.EnableNodePublicIP,
			NodePublicIPPrefixID: agentPool.NodePublicIPPrefixID,
			ScaleSetPriority:     (*string)(agentPool.ScaleSetPriority),
		},
	}

	return worker
}

func (workers *Workers) Convert() (*api.Workers, error) {
	return &api.Workers{
		Defaults: api.DefaultsWorker{
			AzureDefaultWorker: Defaults(),
		},
		WorkersSpec: api.WorkersSpec{
			AzureWorkers: workers.Workers(),
		},
	}, nil
}

func NewAzureWorkers(subscriptionId, resourceGroup string, cluster *armcontainerservice.ManagedCluster) *Workers {
	return &Workers{
		Cluster:        cluster,
		ResourceGroup:  resourceGroup,
		SubscriptionID: subscriptionId,
	}
}
