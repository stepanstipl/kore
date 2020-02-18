## **Kubernetes Deployment**

The following contains a rough guideline as to how kore is deployed into a kubernetes clusters.

```shell
kubectl apply -f deploy/kube/namespaces.yml

# create the secret for the mysql database
kubectl -n kore create secret generic db --from-literal=MYSQL_ROOT_PASSWORD=<CHANGE_ME>

# Open the example secrets folder in deploy/kube/secrets and fill in the details
# Create the various secret consumed by the UI and API
kubectl -n kore create secret generic auth --from-env-file=deploy/kube/secrets/openid.example.env
kubectl -n kore create secret generic kore --from-env-file=deploy/kube/secrets/kore.example.env
kubectl -n kore create secret generic portal --from-env-file=deploy/kube/secrets/portal.example.env

# Create the certificate authority which is consumed by the API - naturely you should
# generat your own self-signed ca's for this.
kubectl -n kore create secret tls ca --cert=hack/ca/ca.pem --key=hack/ca/ca-key.pem

# deploy the api, ui and dependencies
kubectl -n kore apply -f deploy/kube/controlplane.yml
kubectl -n kore apply -f deploy/kube/database.yml
kubectl -n kore apply -f deploy/kube/kore.yml
kubectl -n kore apply -f deploy/kube/portal.yml
```
