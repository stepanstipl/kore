# How to develop locally using kind

Kore depends on and manages it's dependencies in Kubernetes.

This guide enables the use of Kind when testing kore localy.

1. Install kind:

    ```
    GO111MODULE="on" go get sigs.k8s.io/kind@v0.8.1
    ```

1. Create the `charts/my_values.yaml` file:

    set IDP values e.g. see [here](local-quick-start.md#configure-auth0)
    ```
    export KORE_IDP_CLIENT_ID=<your client id>
    export KORE_IDP_CLIENT_SECRET=<your client secret>
    export KORE_IDP_SERVER_URL=<your openid server url>
    ```
    create the helm values file:
    ```
    cat >> ./charts/my_values.yaml << EOF
    ---
    idp:
      client_id: ${KORE_IDP_CLIENT_ID:?'missing var'}
      client_secret: ${KORE_IDP_CLIENT_SECRET:?'missing var'}
      server_url: ${KORE_IDP_SERVER_URL:?'missing var'}
    api:
      endpoint:
        detect: false
      serviceType: NodePort
      hostPort: 10080
      version: latest
      replicas: 1
      feature_gates: []
      version: dev
    ui:
      endpoint:
        detect: false
      serviceType: NodePort
      hostPort: 3000
      version: latest
      replicas: 1
      feature_gates: []
      version: dev
     mysql:
       pvc:
         size: 1Gi
    EOF
    ```

## Local development

Note: Only the API server runs in the Kubernetes cluster. You have to run the UI on your host, as it's really slow in Kind.

1. Build the dev Docker image and load it into the Kind cluster:

    ```
    make kind-dev
    ```

2. Re-build and reload api server changes:

   ```
   make kind-apiserver-reload
   ```

3. Tail api server logs simply

   ```
   make kind-apiserver-logs
   ```
