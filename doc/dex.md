# Dex

See https://github.com/dexidp/dex

DEX provides an identity brokering service by proving a generic OIDC client and upstream "connectors".

Internally DEX stores it's configuration either statically or using CRD's

Dex should be configured through the hub IDP and IDPClient types not directly but the notes below descibe how to configure DEX using the DEX CRD's direct for troubleshooting only.

### Github with DEX CRD
GitHub Authentication can be configured with a CLI client for docker-compose by running the command below.

See "OAuth Apps" under "Developer settings" in GitHub for your org and prepare the statement below with the correct values:

```
export BASE64_GITHUB_OAUTH=$( echo -n '{
  "clientID": "REPLACE ME",
  "clientSecret": "REPLACE ME",
  "redirectURI": "http://127.0.0.1:5556/callback",
  "orgs":
    [{
      "name": "REPLACE ME"
    }]
  }' | base64 )

eval "echo \"$(cat ./hack/setup/dex/example-dex-github-connector.yml)\"" | KUBECONFIG="none" kubectl apply -f -
```

An example OIDC client can be configured to use dex:
```
KUBECONFIG="none" kubectl apply -f ./hack/setup/dex/example-dex-client.yml
```

#### DEX Sample Client

After the CRD's are loaded run the following to start the [dex OIDC sample client](https://github.com/dexidp/dex/blob/master/Documentation/getting-started.md#running-a-client):

```
./bin/example-app --issuer=http://127.0.0.1:5556/
```

Next navigate to the browser and visit the OIDC client http://localhost:5555/ and fill out the form with the client id

### Dex - Issues

DEX currently only supports the API OR static config as the 
source of truth. The backend store is used for recovery on startup.

In addition OAuth2Clients use encoded values for ID's which are decoded when retrieved so we need to encode values using the the right Golang to work with the current CRD's.
