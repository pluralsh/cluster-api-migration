package worker

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2022-03-01/containerservice"
	"github.com/pluralsh/cluster-api-migration/pkg/api"
)

// getDefaultAzureWorkers returns nullified list of default Azure workers.
// It is done as we don't want to create them during migration.
// This has to be kept in sync with bootstrap/helm/cluster-api-cluster/values.yaml.
func getDefaultAzureWorkers() api.AzureWorkers {
	return api.AzureWorkers{
		"lsod":   nil,
		"lsspot": nil,
		"msod":   nil,
		"msspot": nil,
		"ssod":   nil,
		"ssspot": nil,
	}
}

type Workers struct {
	Cluster        *containerservice.ManagedCluster
	ResourceGroup  string
	SubscriptionID string
}

func (workers *Workers) Workers() *api.AzureWorkers {
	result := getDefaultAzureWorkers()

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

func Worker(agentPool containerservice.ManagedClusterAgentPoolProfile) *api.AzureWorker {
	return &api.AzureWorker{
		Replicas:          int(*agentPool.Count),
		KubernetesVersion: nil,
		Annotations:       map[string]string{},
		IsMultiAZ:         true, // default to true so that the availability zones we discovered are used
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
			ScaleDownMode:        (*string)(&agentPool.ScaleDownMode),
			SpotMaxPrice:         agentPool.SpotMaxPrice,
			Scaling: &api.ManagedMachinePoolScaling{
				MinSize: *agentPool.MinCount,
				MaxSize: *agentPool.MaxCount,
			},
		},
	}
}

func (workers *Workers) Convert() (*api.Workers, error) {
	return &api.Workers{
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
