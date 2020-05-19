#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

createkindcluster() {
cat <<EOF | kind create cluster --name kore --wait 1m --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  image: kindest/node:v1.15.11@sha256:6cc31f3533deb138792db2c7d1ffc36f7456a06f1db5556ad3b6927641016f50
EOF
}

if ! kind get clusters | grep "kore" ; then
    createkindcluster
fi

kubectl config use-context kind-kore
export KORE_ENABLE_MANAGED_DEPS=true
#export KUBE_API_SERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
export KUBE_CONFIG_FILE=${KUBECONFIG:-"${HOME}/.kube/config"}
./hack/bin/run-api.sh
