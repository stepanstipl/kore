## **Quick Minikube Deployment**

```shell
minikube start
kubectl config use-context minikube
kubectl apply -f deploy/kube/namespaces.yml

# create the secrets
kubectl -n kore create secret generic dex --from-file=config.yml=deploy/kube/example.dex.yml
kubectl -n kore create secret generic db --from-literal=MYSQL_ROOT_PASSWORD=pass
kubectl -n kore create secret generic kore --from-env-file=deploy/kube/example.kore.env
kubectl -n kore create secret generic portal --from-env-file=deploy/kube/example.portal.env
kubectl -n kore create secret tls certs --cert=hack/ca/ca.pem --key=hack/ca/ca-key.pem

# deploy the api, ui and dependencies
kubectl -n kore apply -f deploy/kube/dex.yml
kubectl -n kore apply -f deploy/kube/api.yml
kubectl -n kore apply -f deploy/kube/portal.yml
```