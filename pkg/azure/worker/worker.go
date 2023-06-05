package worker

import (
	"strings"

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
		result[*agentPool.Name] = Worker(agentPool)
	}

	return &result
}

// AzureTaint returns Azure worker taints mapped from key=value:NoSchedule
// form to taint objects used in CAPI.
func AzureTaint(agentPool *armcontainerservice.ManagedClusterAgentPoolProfile) []api.AzureTaint {
	taints := []api.AzureTaint{}

	for _, taint := range agentPool.NodeTaints {
		effectSplit := strings.Split(*taint, ":")
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

	return taints
}

func Worker(agentPool *armcontainerservice.ManagedClusterAgentPoolProfile) api.AzureWorker {
	worker := api.AzureWorker{
		Replicas:    int(*agentPool.Count),
		Annotations: map[string]string{},
		Spec: api.AzureWorkerSpec{
			AdditionalTags:       agentPool.Tags,
			Mode:                 string(*agentPool.Mode),
			SKU:                  *agentPool.VMSize,
			OSDiskSizeGB:         agentPool.OSDiskSizeGB,
			AvailabilityZones:    agentPool.AvailabilityZones,
			NodeLabels:           agentPool.NodeLabels,
			Taints:               AzureTaint(agentPool),
			MaxPods:              agentPool.MaxPods,
			OsDiskType:           (*string)(agentPool.OSDiskType),
			OSType:               (*string)(agentPool.OSType),
			EnableNodePublicIP:   agentPool.EnableNodePublicIP,
			NodePublicIPPrefixID: agentPool.NodePublicIPPrefixID,
			ScaleSetPriority:     (*string)(agentPool.ScaleSetPriority),
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

func NewAzureWorkers(subscriptionId, resourceGroup string, cluster *armcontainerservice.ManagedCluster) *Workers {
	return &Workers{
		Cluster:        cluster,
		ResourceGroup:  resourceGroup,
		SubscriptionID: subscriptionId,
	}
}
