#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

HUB_APIS=$(find pkg/apis -name "v*" -type d | sed -e 's/pkg\/apis\///' -e 's/\//:/' | sort | tr '\n' ' ')
if [[ -z "${HUB_APIS}" ]]; then 
  echo "[error] unable to find kore apis"
  exit 1
fi

/usr/bin/env bash vendor/k8s.io/code-generator/generate-groups.sh deepcopy  \
  github.com/appvia/kore/pkg/client \
  github.com/appvia/kore/pkg/apis \
  "${HUB_APIS}" \
  --go-header-file hack/boilerplate.go.txt
