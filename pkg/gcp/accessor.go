package gcp

import (
	"context"
	"fmt"
	"log"

	cloudcontainer "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"google.golang.org/api/option"
	"k8s.io/client-go/pkg/version"

	"github.com/pluralsh/cluster-api-migration/pkg/api"
	"github.com/pluralsh/cluster-api-migration/pkg/gcp/cluster"
)

type ClusterAccessor struct {
	configuration *api.GCPConfiguration
	ctx           context.Context
	client        *cloudcontainer.ClusterManagerClient
}

func (this *ClusterAccessor) init() api.ClusterAccessor {
	client, err := cloudcontainer.NewClusterManagerClient(
		this.ctx,
		this.defaultClientOptions(this.configuration.Credentials)...,
	)

	if err != nil {
		log.Fatal(err)
	}

	this.client = client
	return this
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
	c, err := this.client.GetCluster(this.ctx, &containerpb.GetClusterRequest{
		Name: this.clusterName(this.configuration.Project, this.configuration.Region, this.configuration.Name),
	})

	if err != nil {
		log.Fatal(err)
	}

	gcpCluster := cluster.NewGCPCluster(this.configuration.Project, c)
	return gcpCluster.Convert()
}

func (this *ClusterAccessor) GetWorkers() *api.Workers {
	return &api.Workers{
		Defaults:    api.DefaultsWorker{},
		WorkersSpec: api.WorkersSpec{},
	}
}
