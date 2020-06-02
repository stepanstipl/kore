# How to install Kore in Kind using the Helm chart

1. Install kind:

    ```
    GO111MODULE="on" go get sigs.k8s.io/kind@v0.8.1
    ```

1. Create a Kind cluster:

    ```
    cat <<EOF | kind create cluster --name kore --config=-
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    nodes:
      - role: control-plane
        image: kindest/node:v1.15.11@sha256:6cc31f3533deb138792db2c7d1ffc36f7456a06f1db5556ad3b6927641016f50
        extraPortMappings:
          - containerPort: 3000
            hostPort: 3000
            protocol: TCP
          - containerPort: 10080
            hostPort: 10080
            protocol: TCP
        extraMounts:
          - hostPath: ${GOPATH}/src/github.com/appvia/kore
            containerPath: /go/src/github.com/appvia/kore
    EOF
    ```

1. Make sure to use the kind context for kore:

    ```
    kubectl config use-context kind-kore
    ```

1. Create the kore namespace

    ```
    kubectl create ns kore
    ```

1. Create the `charts/my_values.yaml` file:

    ```
    ---
    idp:
      client_id: [...]
      client_secret: [...]
      server_url: [...]
    api:
      endpoint:
        detect: false
      serviceType: NodePort
      hostPort: 10080
      version: latest
      replicas: 1
      feature_gates: []
    ui:
      endpoint:
        detect: false
      serviceType: NodePort
      hostPort: 3000
      version: latest
      replicas: 1
      feature_gates: []
     mysql:
       pvc:
         size: 1Gi
    ```

1. Install the Helm chart

    Ensure you have helm v3 installed: https://github.com/helm/helm/releases

    ```
    helm install --namespace kore kore ./charts/kore --wait -f ./charts/my_values.yaml
    ```

1. Navigate to http://localhost:3000 or run `kore login`

## Local development

1. Build the dev Docker image and load it into the Kind cluster:

    ```
    make kind-image-dev
    ```

1. Update `charts/my_values.yaml` and set `api.version` and/or `ui.version` to `dev`

### Useful commands

#### Restart the API server

   ```
   make kind-apiserver-reload
   ```

#### Tail the API server logs

   ```
   make kind-apiserver-logs
   ```

#### Restart the UI

   ```
   make kind-ui-reload
   ```

#### Tail the UI logs

   ```
   make kind-ui-logs
   ```
