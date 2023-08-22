package gcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/version"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp/cluster"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp/worker"
)

type ClusterAccessor struct {
	configuration    *api.GCPConfiguration
	ctx              context.Context
	clusterClient    *container.ClusterManagerClient
	computeClient    *compute.Service
	kubernetesClient *kubernetes.Clientset
}

func (this *ClusterAccessor) AddClusterTags(Tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) AddMachinePollsTags(Tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) AddVirtualNetworkTags(tags map[string]string) error {
	return nil
}

func (this *ClusterAccessor) init() (api.ClusterAccessor, error) {
	err := this.initContainerClient()
	if err != nil {
		return nil, err
	}

	err = this.initComputeClient()
	if err != nil {
		return nil, err
	}

	err = this.initKubernetesClient()
	return this, err
}

func (this *ClusterAccessor) initContainerClient() error {
	client, err := container.NewClusterManagerClient(
		this.ctx,
		this.defaultClientOptions()...,
	)

	if err != nil {
		return err
	}

	this.clusterClient = client
	return nil
}

func (this *ClusterAccessor) initComputeClient() error {
	client, err := compute.NewService(
		this.ctx,
		this.defaultClientOptions()...,
	)

	if err != nil {
		return err
	}

	this.computeClient = client
	return nil
}

func (this *ClusterAccessor) initKubernetesClient() error {
	config, err := clientcmd.BuildConfigFromFlags("", this.configuration.KubeconfigPath)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	this.kubernetesClient = clientset
	return nil
}

func (this *ClusterAccessor) defaultClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithUserAgent(fmt.Sprintf("gcp.cluster.x-k8s.io/%s", version.Get())),
	}
}

func (this *ClusterAccessor) clusterName(project, region, name string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
		project,
		region,
		name,
	)
}

func (this *ClusterAccessor) getClusterOrDie() *containerpb.Cluster {
	cluster, err := this.clusterClient.GetCluster(this.ctx, &containerpb.GetClusterRequest{
		Name: this.clusterName(this.configuration.Project, this.configuration.Region, this.configuration.Name),
	})

	if err != nil {
		log.Fatal(err)
	}

	return cluster
}

func (this *ClusterAccessor) getNetworkOrDie(name string) *compute.Network {
	req := this.computeClient.Networks.Get(this.configuration.Project, name)
	network, err := req.Do()

	if err != nil {
		log.Fatal(err)
	}

	return network
}

func (this *ClusterAccessor) getSubnetworksOrDie(network string) []*compute.Subnetwork {
	result := make([]*compute.Subnetwork, 0)
	r := this.computeClient.Subnetworks.List(this.configuration.Project, this.configuration.Region)
	if err := r.Pages(this.ctx, func(page *compute.SubnetworkList) error {
		for _, subnetwork := range page.Items {
			if strings.HasSuffix(subnetwork.Network, network) {
				result = append(result, subnetwork)
			}
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return result
}

func (this *ClusterAccessor) getNodesOrDie() *corev1.NodeList {
	nodes, err := this.kubernetesClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return nodes
}

func (this *ClusterAccessor) GetCluster() (*api.Cluster, error) {
	c := this.getClusterOrDie()
	network := this.getNetworkOrDie(c.Network)
	subnetworks := this.getSubnetworksOrDie(network.Name)
	gcpCluster := cluster.NewGCPCluster(this.configuration.Project, c, network, subnetworks)

	return gcpCluster.Convert(), nil
}

func (this *ClusterAccessor) GetWorkers() (*api.Workers, error) {
	cluster := this.getClusterOrDie()
	nodes := this.getNodesOrDie()
	workers := worker.NewGCPWorkers(cluster, nodes)

	return workers.Convert(), nil
}
