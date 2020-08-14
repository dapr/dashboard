# ClusterRole

If Dashboard needs to access other Kubernetes resources, the deployment-reader clusterrole needs to be updated for testing:
```bash
kubectl delete clusterrole deployment-reader
```

Apply your changes to the new clusterrole file:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: deployment-reader
rules:
- apiGroups: ["", "dapr.io", "apps", "extensions"]
  resources: ["deployments", ..., "<new-resource-type>"]
  verbs: ["get", "list", ..., "<new-verb-type>"]
```

Apply the new clusterrole to the cluster:
```bash
kubectl apply -f ./test_clusterrole
```

These changes must also be made in the Helm charts in the [dapr/dapr](https://github.com/dapr/dapr/tree/master/charts/dapr/charts/dapr_dashboard) repository for deployment in the same fashion.