## Updating application controller

The project only maintain a dev image at this time.

### Build the conatainer
The image referred to was generated from the source e.g.:
```
APP_VERSION=v0.8.2

mkdir -p ${GOPATH}/src/github.com/kubernetes-sigs
cd ${GOPATH}/src/github.com/kubernetes-sigs
git clone git@github.com:kubernetes-sigs/application.git
cd application

git checkout ${APP_VERSION}

docker build -t quay.io/appvia/application-controller:${APP_VERSION} .
docker push quay.io/appvia/application-controller:${APP_VERSION}
```

### Update Manifests
The manifests were taken from here (with the image updated as below):
```
APP_VERSION=v0.8.2
curl -sSL https://raw.githubusercontent.com/kubernetes-sigs/application/${APP_VERSION}/deploy/kube-app-manager-aio.yaml > ./pkg/kore/assets/applications/kube-app-manager/kube-app-manager-aio.yaml

// edit the image to match our container:
vi ./pkg/kore/assets/applications/kube-app-manager/kube-app-manager-aio.yaml
```
