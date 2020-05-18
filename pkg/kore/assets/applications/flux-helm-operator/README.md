# Update manifests

1. Download the manifests from the fluxcd/helm-operator Github repo:

    ```
    (
    cd pkg/kore/assets/applications/flux-helm-operator
    FLUX_VERSION=v1.0.1
    
    for file in crds.yaml deployment.yaml namespace.yaml rbac.yaml; do
      curl -sSL "https://raw.githubusercontent.com/fluxcd/helm-operator/${FLUX_VERSION}/deploy/$file" -o "$file"
    done
    )
    ```

1. Enabled Helm 3 support in `deployment.yaml`

    ```
    - --enabled-helm-versions=v3
    ```

1. Disable Tiller in `deployment.yaml`

    ```
    #- --tiller-namespace=kube-system
    ```

1. Add the app label to `deployment.yaml`:

    ```
    metadata:
      labels:
        app.kubernetes.io/name: flux-helm-operator
    ```
