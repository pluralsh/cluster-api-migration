# cluster-api-migration

## Testing

To test migrations use following repos/branches:

- <https://github.com/pluralsh/plural-cli/tree/bootstrap>
- <https://github.com/pluralsh/plural-artifacts/tree/cluster-api-cluster>

### Azure

Setup Azure cluster using old method:

```sh
plural init
plural bundle install bootstrap azure-k8s
plural build
plural deploy
```

Once AKS is up and running you can start the migration by generating `values.yaml` using this repo and move it to artifacts repo:

```sh
cp $WORKSPACE/values.yaml $WORKSPACE/plural-artifacts/bootstrap/helm/cluster-api-cluster/
```

Set following tags on AKS:

- `sigs.k8s.io_cluster-api-provider-azure_cluster_aaa` : `owned`
- `sigs.k8s.io_cluster-api-provider-azure_role` : `common`

Disable `azure-identity` by setting `azure-identity.enabled` to `false` in `aaa/bootstrap/helm/bootstrap/default-values.yaml` (`aaa` is the name of installation repo). Ensure that it wasn't changed after each build.

Install new recipe:

```sh
plural bundle install bootstrap azure-cluster-api
plural build --cluster-api --force
plural link helm bootstrap --name bootstrap-operator --path $WORKSPACE/plural-artifacts/bootstrap/helm/bootstrap-operator/
plural link helm bootstrap --name cluster-api-cluster --path $WORKSPACE/plural-artifacts/bootstrap/helm/cluster-api-cluster/
```

Rebuild:

```sh
plural build --cluster-api --force
```

Deploy CRDs and Helm chart:
```sh
cd bootstrap
plural workspace crds bootstrap
plural workspace helm bootstrap --skip cluster-api-cluster
sleep 120
plural workspace helm bootstrap
```
