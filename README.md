# cluster-api-migration

## Testing

To test migrations use following repos/branches:

- https://github.com/pluralsh/plural-cli/tree/bootstrap
- https://github.com/pluralsh/plural-artifacts/tree/cluster-api-cluster

### Azure

Setup Azure cluster using old method:

```sh
plural install
plural bundle install bootstrap azure-k8s
plural build
plural deploy
```

Once AKS is up and running you can start the migration.

Generate new `values.yaml` using this repo and move it to artifacts repo:

```sh
cp $WORKSPACE/values.yaml $WORKSPACE/plural-artifacts/bootstrap/helm/cluster-api-cluster/
```

Then you can install new recipe:

```sh
plural bundle install bootstrap azure-cluster-api
plural build --cluster-api
plural link helm bootstrap --name cluster-api-cluster --path $WORKSPACE/plural-artifacts/bootstrap/helm/cluster-api-cluster/
```

Then go to your installation repo (in this case `aaa`) rebuild it, deploy CRDs and Helm chart:

```sh
cd $WORKSPACE/aaa
plural build --cluster-api --force
cd $WORKSPACE/aaa/bootstrap
plural workspace crds bootstrap
sleep 30
plural workspace helm bootstrap --skip cluster-api-cluster
sleep 60
plural workspace helm bootstrap
```

At the moment it is required to set following AKS tags manually:

- `sigs.k8s.io_cluster-api-provider-azure_cluster_aaa` : `owned`
- `sigs.k8s.io_cluster-api-provider-azure_role` : `common`
