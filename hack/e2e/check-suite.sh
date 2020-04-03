#!/bin/bash
#
# Copyright 2020 Appvia Ltd <info@appvia.io>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

source hack/e2e/environment.sh || exit 1

# These's are using across the checks
export CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"
export E2E_DIR="e2eci"
export KORE_IDP_SERVER_URL=${KORE_IDP_SERVER_URL:-"unknown"}
export KORE_IDP_CLIENT_ID=${KORE_IDP_CLIENT_ID:-"unknown"}
export KORE_API_URL=${KORE_API_URL_QA:-"http://127.0.0.1:10080"}
export KORE_ID_TOKEN=${KORE_ID_TOKEN_QA:-"unknown"}
export KORE_PROFILE="local"
#export KORE_USERNAME="${KORE_USERNAME:-"kore-ci@appvia.io"}"

ENABLE_CONFORMANCE=${ENABLE_CONFORMANCE:-false}
ENABLE_UNIT_TESTS=${ENABLE_UNIT_TESTS:-true}
MAX_RETRIES=${1:-60}
PATH=$PATH:${GOPATH}/bin
KORE_CONFIG=${HOME}/.korectl/config
RETRIES=0
WAIT_TIME=20

mkdir -p ${GOPATH}/bin

usage() {
  cat <<EOF
  Usage: $(basename $0)
  --enable-conformance   : run the kubernetes conformance check suite
  --enable-unit-tests    : run the bats unit tests
  -h|--help              : display this usage menu
EOF
  if [[ -n $@ ]]; then
    echo "[error] $@"
    exit 1
  fi
  exit 0
}

create-korectl-config() {
  [[ -f ${KORE_CONFIG} ]] && return

  announce "Generating a korectl configuration: ${KORE_CONFIG}"
  cat << EOF > ${KORE_CONFIG}
current-profile: local
profiles:
  local:
    server: local
    user: local
servers:
  local:
    server: ${KORE_API_URL}
users:
  local:
    oidc:
      client-id: ${KORE_IDP_CLIENT_ID}
      id-token: ${KORE_ID_TOKEN}
      idp-issuer-url: ${KORE_IDP_SERVER_URL}
EOF
}

wait-kubeapi-readiness() {
  announce "waiting for kube-apiserver readiness ..."

  while true ; do
    if kubectl version &> /dev/null; then
      if kubectl get namespace >/dev/null; then
        break
      fi
    fi

    RETRIES=$((RETRIES + 1))
    if [[ ${RETRIES} -eq ${MAX_RETRIES} ]]; then
      error "max timeout reached. kube-apiserver not ready ..."
      kubectl --namespace=kube-system get pods || true
      exit 1
    else
      announce "attempt #${RETRIES} of #${MAX_RETRIES}: dns not yet available, sleeping for ${WAIT_TIME} seconds..."
      sleep ${WAIT_TIME}
    fi
  done

  return 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
  --enable-conformance) ENABLE_CONFORMANCE="true"; shift 1; ;;
  --enable-unit-tests)  ENABLE_UNIT_TESTS="true";  shift 1; ;;
  -h|--help)            usage;                            ;;
  *)                                             shift 1; ;;
  esac
done

## @step: check the api server is running and available
#if ! wait-kubeapi-readiness; then
#  error "kube apiserver is not available after multiple attempts"
#  exit 1
#fi

if [[ "${ENABLE_UNIT_TESTS}" == "true" ]]; then
  if [[ ! -f "${E2E_DIR}/gke-credentials.yml" ]]; then
    error "you need to have the ${E2E_DIR}/gke-credentials.yml containing real credentials for gke"
    error "you can copy an template from: examples/gke.yml"
    exit 1
  fi

  if ! create-korectl-config; then
    error "failed to generate a client configuration for cli"
    exit 1
  fi

  hack/e2e/check-suite-units.sh || exit 1
else
  announce "skipping the unit tests suite"
fi

if [[ "${ENABLE_CONFORMANCE}" == "true" ]]; then
  hack/e2e/check-suite-conformance.sh || exit 1
else
  announce "skipping the kubernetes e2e conformance suite"
fi
