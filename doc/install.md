## Install

## GKE Quickstart

### 1. Create a GKE cluster and configure access

You can find instructions for this [here](https://cloud.google.com/kubernetes-engine/docs/quickstart). Follow them to the end of the 'Get authentication credentials for the cluster' step

### 2. Initial IDP Configuration
Configure Auth0 by following the docs [here](local-quick-start.md#configure-auth0) and export the obtained values. You do not need to configure 'Allowed Callback URLs' at this stage

```
export KORE_IDP_CLIENT_ID=<your client id>
export KORE_IDP_CLIENT_SECRET=<your client secret>
export KORE_IDP_SERVER_URL=<your openid server url>
```

### 3. Create helm configuration values
Create a file with the correct values
```
# Create a localy ignored file:
cat >> ./charts/my_values.yaml << EOF
idp:
  client_id: $KORE_IDP_CLIENT_ID
  client_secret: $KORE_IDP_CLIENT_SECRET
  server_url: $KORE_IDP_SERVER_URL
api:
  endpoint:
    detect: true
  serviceType: LoadBalancer
ui:
  endpoint:
    detect: true
  serviceType: LoadBalancer
EOF
```

In the above commands, the services are of the LoadBalancer type and endpoint detection is enabled. Endpoint detection will allow the LoadBalancer assigned ip addresses to be automatically configured by the kore api and ui

### 4. Install Kore

1. Create the namespace for kore installation

`kubectl create ns kore`

2. Install helm using your configured values

`helm install --namespace kore kore ./charts/kore --wait -f ./charts/my_values.yaml`

The helm installation will wait for LoadBalancers to be configured before completing

### 5. Configure Auth0

You will need the loadbalancer endpoints for the kore-apiserver. These can be found by running:

`kubectl -n kore get services`

Add the following values to the `Allowed Callback Urls` section of Auth0:
```
http://<apiserver external ip>:10080/oauth/callback
http://<portal external ip>:3000/auth/callback
```

### 6. Log In

Visit the ui at `http://<portal ip>:3000` or configure kore using `kore login <new profile name> -a http://<portal ip>:10080`

