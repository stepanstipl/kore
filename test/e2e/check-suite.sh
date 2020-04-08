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

source test/e2e/environment.sh || exit 1

# These are being used across the checks
export CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"
export E2E_DIR="e2eci"
export KORE_API_URL=${KORE_API_PUBLIC_URL_QA:-"http://127.0.0.1:10080"}
export KORE_IDP_CLIENT_ID=${KORE_IDP_CLIENT_ID:-"unknown"}
export KORE_IDP_SERVER_URL=${KORE_IDP_SERVER_URL:-"unknown"}
export KORE_ID_TOKEN=${KORE_ID_TOKEN_QA:-"unknown"}
export KORE_PROFILE="local"
export TEAM="e2e"

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
  --enable-e2e-user      : indicates we should provision a test user
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

  mkdir -p $(dirname ${KORE_CONFIG}) || {
    error "unable to create the client configuration directory";
    exit 1;
  }

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
      authorize-url: ${KORE_IDP_SERVER_URL}
      client-id: ${KORE_IDP_CLIENT_ID}
      id-token: ${KORE_ID_TOKEN}
EOF
}

enable-admin-user() {
  announce "Provisioning a admin e2e user"
  cat <<EOF >/tmp/e2e.user
{
  "apiVersion": "org.kore.appvia.io/v1",
  "kind": "User",
  "metadata": {
    "name": "${ADMIN_USER}"
  },
  "spec": {
    "username": "${ADMIN_USER}",
    "email": "${ADMIN_USER}"
  }
}
EOF

  if ! curl -sL -X PUT \
    --retry 3 \
    --header "Content-Type: application/json" \
    --header "Authorization: Bearer ${KORE_ADMIN_TOKEN}" \
    ${KORE_API_URL}/api/v1alpha1/users/${ADMIN_USER} -d @/tmp/e2e.user; then
    error "trying to provision admin user for e2e"
    exit 1
  fi

  if ! curl -sL -X PUT \
    --retry 3 \
    --header "Content-Type: application/json" \
    --header "Authorization: Bearer ${KORE_ADMIN_TOKEN}" \
    ${KORE_API_URL}/api/v1alpha1/teams/kore-admin/members/${ADMIN_USER}; then
    error "trying to provision team membership"
    exit 1
  fi

  return 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
  --enable-conformance) ENABLE_CONFORMANCE="true"; shift 1; ;;
  --enable-unit-tests)  ENABLE_UNIT_TESTS="true";  shift 1; ;;
  --enable-e2e-user)    ADMIN_USER=$2;             shift 2; ;;
  -h|--help)            usage;                              ;;
  *)                                               shift 1; ;;
  esac
done

if [[ -n "${ADMIN_USER}" ]]; then
  enable-admin-user
fi

if [[ "${ENABLE_UNIT_TESTS}" == "true" ]]; then
  if [[ -n "${GKE_SA_QA}" ]]; then
    mkdir -p ${E2E_DIR}
    echo -n ${GKE_SA_QA} | base64 -d > ${E2E_DIR}/gke-credentials.yml 2>/dev/null
  fi

  if [[ ! -f "${E2E_DIR}/gke-credentials.yml" ]]; then
    error "you need to have the ${E2E_DIR}/gke-credentials.yml containing real credentials for gke"
    error "you can copy an template from: examples/gke.yml"
    exit 1
  fi

  if ! create-korectl-config; then
    error "failed to generate a client configuration for cli"
    exit 1
  fi

  test/e2e/check-suite-units.sh || exit 1
else
  announce "skipping the unit tests suite"
fi

if [[ "${ENABLE_CONFORMANCE}" == "true" ]]; then
  test/e2e/check-suite-conformance.sh || exit 1
else
  announce "skipping the kubernetes e2e conformance suite"
fi
