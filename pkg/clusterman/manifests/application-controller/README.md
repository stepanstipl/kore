See issue https://github.com/kubernetes-sigs/application/issues/141

The manifests were taken from here:
```
git clone git@github.com:appvia/application.git
git checkout v0.8.1-patched-image-version
kubectl kustomize ./config/ > application-all.yaml
```

The image refered to was generated from the source above with:
```
docker build -t quay.io/appvia/application-controller:v0.8.1 .
docker push quay.io/appvia/application-controller:v0.8.1
```
