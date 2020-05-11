Update manifests

```
(
cd pkg/clusterappman/manifests/flux
FLUX_VESION=v1.0.1

for file in crds.yaml deployment.yaml namespace.yaml rbac.yaml ; do
  curl -sSL https://raw.githubusercontent.com/fluxcd/helm-operator/${VERSION}/deploy/$file -o $file
done
)
```

TODO add kustomize to ensure we have the following on the deployment:
```
  labels:
    app.kubernetes.io/name: flux-helm-operator

```