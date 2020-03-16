# How to install Kore in Kind using the Helm chart

1. Install kind:

    ```
    GO111MODULE="on" go get sigs.k8s.io/kind@v0.7.0
    ```

1. Create a Kind cluster:

    ```
    cat <<EOF | kind create cluster --name kore --config=-
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    nodes:
    - role: control-plane
      image: kindest/node:v1.14.10@sha256:81ae5a3237c779efc4dda43cc81c696f88a194abcc4f8fa34f86cf674aa14977
      extraPortMappings:
      - containerPort: 3000
        hostPort: 3000
        protocol: TCP
      - containerPort: 10080
        hostPort: 10080
        protocol: TCP
    EOF
    ```

1. Make sure to use the kind context for korectl:

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
      version: v0.0.16
    ui:
      endpoint:
        detect: false
      serviceType: NodePort
      hostPort: 3000
      version: v0.0.6
    ```

1. Install the Helm chart

    ```
    helm install --namespace kore kore ./charts/kore --wait -f ./charts/my_values.yaml
    ```

1. Navigate to http://localhost:3000 or run `korectl login`
