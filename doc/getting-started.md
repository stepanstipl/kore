# **Development**

In order to setup a local environment

## Pre-requisites

* Install and start [Docker](https://www.docker.com/products/docker-desktop)
* docker-compose (which comes part of Docker)
* Make sure you use Go 1.13+

## Build and Run

### Build

If you are building from master and not a known version, you will need to make sure the pre-requisites are met.

The following commands will build both the apiserver and the CLI:

`make compose`

This will create dex configuration for the authentication, placing the DNS name inside of the configuration. The output of this is inside: `./hack/setup/dex/config.yaml`

The kube-controller-manager, etcd, mysql and the kube-apiserver will also then be started, which are the core dependencies for the kore apiserver. For information on those they can be found in:

`hack/compose/kube.yml`

The dex operator is then added, with a path to the hack directory and the config.yaml for dex sourced.

### Run

Ideally you can set Environment variables to the kore-apiserver:

```
export DATABASE_URL="$USER:$PASS@tcp(:3306)/hub?parseTime=true"
export HUB_ADMIN_TOKEN=<ADMIN_TOKEN>
export HUB_CERTIFICATE_AUTHORITY_KEY=hack/ca/ca-key.pem
export HUB_CERTIFICATE_AUTHORITY=hack/ca/ca.pem
export HUB_CLIENT_ID=<OPENID CLIENT ID>
export HUB_CLIENT_SECRET=<OPENID SECRET>
export HUB_DISCOVERY_URL=<OPENID DISCOVERY URL>
export HUB_AUTHENTICATION_PLUGINS=openid,admintoken
```

You can override everything in here with what you need. Alternatively you can run:

`./bin/kore-apiserver-h` 

For a list of the options

Once this start it will create a set of bootstrap namespaces, default teams etc. From that point on, you can then begin to use: `./bin/korectl` to talk to the API server to begin provisioning resources

### Swagger UI

You can view the swagger at `http://127.0.0.1:10080/swagger.json`. Note if you want to see the pretty swagger UI, can you download the swagger-ui from https://github.com/swagger-api/swagger-ui/. Grab the `dist` folder inside the repo and move to the base swagger-ui/ in this repo. You can then open: http://127.0.0.1:10080/apidocs/?url=http://localhost:10080/swagger.json

### Demo

To run a demo of the hub simply type: `make demo` in the base of the repo. If you want to ensure this is a fresh install use `make clean`

### Auth

See [configure an IDP](./docs/idp.md)
