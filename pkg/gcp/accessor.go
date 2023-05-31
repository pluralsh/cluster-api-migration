package gcp

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"k8s.io/client-go/pkg/version"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp/cluster"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp/worker"
)

type ClusterAccessor struct {
	configuration *api.GCPConfiguration
	ctx           context.Context
	clusterClient *container.ClusterManagerClient
	computeClient *compute.Service
}

func (this *ClusterAccessor) init() api.ClusterAccessor {
	this.initContainerClientOrDie()
	this.initComputeClientOrDie()

	return this
}

func (this *ClusterAccessor) initContainerClientOrDie() {
	client, err := container.NewClusterManagerClient(
		this.ctx,
		this.defaultClientOptions(this.configuration.Credentials)...,
	)

	if err != nil {
		log.Fatal(err)
	}

	this.clusterClient = client
}

func (this *ClusterAccessor) initComputeClientOrDie() {
	client, err := compute.NewService(
		this.ctx,
		this.defaultClientOptions(this.configuration.Credentials)...,
	)

	if err != nil {
		log.Fatal(err)
	}

	this.computeClient = client
}

func (this *ClusterAccessor) defaultClientOptions(credentials string) []option.ClientOption {
	return []option.ClientOption{
		option.WithUserAgent(fmt.Sprintf("gcp.cluster.x-k8s.io/%s", version.Get())),
		option.WithCredentialsJSON([]byte(credentials)),
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

func (this *ClusterAccessor) getSubnetworkOrDie(name string) *compute.Subnetwork {
	// TODO: Check if we shouldn't search for all subnets and filter by parent network name
	req := this.computeClient.Subnetworks.Get(this.configuration.Project, this.configuration.Region, name)
	subnetwork, err := req.Do()

	if err != nil {
		log.Fatal(err)
	}

	return subnetwork
}

func (this *ClusterAccessor) GetCluster() *api.Cluster {
	c := this.getClusterOrDie()
	network := this.getNetworkOrDie(c.Network)
	subnetwork := this.getSubnetworkOrDie(c.Subnetwork)
	gcpCluster := cluster.NewGCPCluster(this.configuration.Project, c, network, subnetwork)

	return gcpCluster.Convert()
}

func (this *ClusterAccessor) GetWorkers() *api.Workers {
	cluster := this.getClusterOrDie()
	workers := worker.NewGCPWorkers(cluster)

	return workers.Convert()
}
