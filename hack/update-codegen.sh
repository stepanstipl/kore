#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

HUB_APIS=$(find pkg/apis -name "v*" -type d | sed -e 's/pkg\/apis\///' -e 's/\//:/' | sort | tr '\n' ' ')
if [[ -z "${HUB_APIS}" ]]; then 
  echo "[error] unable to find kore apis"
  exit 1
fi

if [[ -x ${GOPATH}/src/k8s.io/code-generator/generate-groups.sh ]]; then
  echo "Installing k8s.io/code-generator/cmd/..."
  export GO111MODULE=off
  # We need to build without modules but with a specified version...
  go get -d k8s.io/code-generator/cmd/...
  (
    cd ${GOPATH}/src/k8s.io/code-generator
    git checkout kubernetes-1.14.1
    go build -o ${GOPATH}/bin/ ./cmd/register-gen/...
  )
fi

${GOPATH}/src/k8s.io/code-generator/generate-groups.sh deepcopy  \
  github.com/appvia/kore/pkg/client \
  github.com/appvia/kore/pkg/apis \
  "${HUB_APIS}" \
  --go-header-file hack/boilerplate.go.txt
