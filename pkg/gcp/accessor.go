package gcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/pluralsh/polly/algorithms"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"k8s.io/client-go/pkg/version"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp/cluster"
	"github.com/pluralsh/cluster-api-migration/pkg/resources"
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

func (this *ClusterAccessor) GetCluster() *api.Cluster {
	c := this.getClusterOrDie()
	network := this.getNetworkOrDie(c.Network)
	subnetworkList := this.getSubnetworkListOrDie(network.Subnetworks)

	resources.NewPrinter(c).PrettyPrint()
	resources.NewPrinter(network).PrettyPrint()
	resources.NewPrinter(subnetworkList).PrettyPrint()

	gcpCluster := cluster.NewGCPCluster(this.configuration.Project, c)
	return gcpCluster.Convert()
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

func (this *ClusterAccessor) getSubnetworkListOrDie(names []string) *compute.SubnetworkList {
	subnetworks := new(compute.SubnetworkList)
	req := this.computeClient.Subnetworks.List(this.configuration.Project, this.configuration.Region)

	if err := req.
		Filter(strings.Join(algorithms.Map(names, func(name string) string {
			return fmt.Sprintf("name=%s", name)
		}), ",")).
		Pages(this.ctx, func(page *compute.SubnetworkList) error {
			subnetworks = page
			return nil
		}); err != nil {
		log.Fatal(err)
	}

	return subnetworks
}

func (this *ClusterAccessor) GetWorkers() *api.Workers {
	return &api.Workers{
		Defaults:    api.DefaultsWorker{},
		WorkersSpec: api.WorkersSpec{},
	}
}
