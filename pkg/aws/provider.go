package aws

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func getCfg() *api.ClusterConfig {
	cfg := api.NewClusterConfig()
	cfg.IAM.WithOIDC = api.Enabled()
	return cfg
}
