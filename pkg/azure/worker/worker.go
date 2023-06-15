package worker

import (
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2022-03-01/containerservice"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"strings"
)

type Workers struct {
	Cluster        *containerservice.ManagedCluster
	ResourceGroup  string
	SubscriptionID string
}

func (workers *Workers) Workers() *api.AzureWorkers {
	result := api.AzureWorkers{}
	for _, agentPool := range *workers.Cluster.AgentPoolProfiles {
		result[*agentPool.Name] = Worker(agentPool)
	}

	return &result
}

func Taints(agentPool containerservice.ManagedClusterAgentPoolProfile) []api.AzureTaint {
	taints := make([]api.AzureTaint, 0)

	if agentPool.NodeTaints != nil {
		for _, taint := range *agentPool.NodeTaints {
			effectSplit := strings.Split(taint, ":")
			if len(effectSplit) >= 2 {
				keyValueSplit := strings.Split(effectSplit[0], "=")
				if len(keyValueSplit) >= 2 {
					taints = append(taints, api.AzureTaint{
						Effect: effectSplit[1],
						Key:    keyValueSplit[0],
						Value:  keyValueSplit[1],
					})
				}
			}
		}
	}

	return taints
}

func NodeLabels(agentPool containerservice.ManagedClusterAgentPoolProfile) map[string]*string {
	labels := make(map[string]*string)
	for key, value := range agentPool.NodeLabels {
		// Node pool label key must not start with kubernetes.azure.com.
		if !strings.HasPrefix(key, "kubernetes.azure.com") {
			labels[key] = value
		}
	}

	return labels
}

func Worker(agentPool containerservice.ManagedClusterAgentPoolProfile) api.AzureWorker {
	worker := api.AzureWorker{
		Replicas:    int(*agentPool.Count),
		Annotations: map[string]string{},
		Spec: api.AzureWorkerSpec{
			AdditionalTags:       agentPool.Tags,
			Mode:                 string(agentPool.Mode),
			SKU:                  *agentPool.VMSize,
			OSDiskSizeGB:         agentPool.OsDiskSizeGB,
			AvailabilityZones:    *agentPool.AvailabilityZones,
			NodeLabels:           NodeLabels(agentPool),
			Taints:               Taints(agentPool),
			MaxPods:              agentPool.MaxPods,
			OsDiskType:           (*string)(&agentPool.OsDiskType),
			OSType:               (*string)(&agentPool.OsType),
			EnableNodePublicIP:   agentPool.EnableNodePublicIP,
			NodePublicIPPrefixID: agentPool.NodePublicIPPrefixID,
			ScaleSetPriority:     (*string)(&agentPool.ScaleSetPriority),
			Scaling: &api.ManagedMachinePoolScaling{
				MinSize: *agentPool.MinCount,
				MaxSize: *agentPool.MaxCount,
			},
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

func NewAzureWorkers(subscriptionId, resourceGroup string, cluster *containerservice.ManagedCluster) *Workers {
	return &Workers{
		Cluster:        cluster,
		ResourceGroup:  resourceGroup,
		SubscriptionID: subscriptionId,
	}
}
