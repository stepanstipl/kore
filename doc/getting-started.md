## **Development**

In order to setup a local environment

* Install and start [Docker](https://www.docker.com/products/docker-desktop)
* Run the following commands:
```shell
$ git clone git@github.com:appvia/kore-apiserver.git
$ cd kore-apiserver
$ make compose
# deploy the crds for the kore
$ git clone git@github.com:appvia/kore-apis.git
$ cd kore-apis
$ KUBECONFIG="none" kubectl apply -f ./deploy
# run the kore
$ cd kore-apiserver
$ export GOPRIVATE=github.com/appvia
$ make
$ bin/kore-apiserver --kube-api-server http://127.0.0.1:8080 --verbose --dex-public-url http://127.0.0.1:5556 --dex-grpc-server 127.0.0.1 --admin-pass xyz
```

### Swagger UI

You can view the swagger at `http://127.0.0.1:10080/swagger.json`. Note if you want to see the pretty swagger UI, can you download the swagger-ui from https://github.com/swagger-api/swagger-ui/. Grab the `dist` folder inside the repo and move to the base swagger-ui/ in this repo. You can then open: http://127.0.0.1:10080/apidocs/?url=http://localhost:10080/swagger.json

### Demo

To run a demo of the kore simply type: `make demo` in the base of the repo. If you want to ensure this is a fresh install use `make clean`

### Auth

See [configure an IDP](./docs/idp.md)