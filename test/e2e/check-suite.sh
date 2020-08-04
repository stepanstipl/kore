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
export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID_QA}
export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY_QA}
export AWS_DEFAULT_REGION="eu-west-2"
export AWS_KORE_PROFILE="kore-e2e"
export CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"
export E2E_DIR="e2eci"
export KORE_API_URL=${KORE_API_PUBLIC_URL_E2E:-"http://localhost:10080"}
export KORE_IDP_CLIENT_ID=${KORE_IDP_CLIENT_ID:-"unknown"}
export KORE_IDP_SERVER_URL=${KORE_IDP_SERVER_URL:-"unknown"}
export KORE_ID_TOKEN=${KORE_ID_TOKEN_QA:-"unknown"}
export KORE_PROFILE="local"
export TEAM="${TEAM:-"e2e"}"

ENABLE_CONFORMANCE=${ENABLE_CONFORMANCE:-false}
ENABLE_EKS_E2E=${ENABLE_EKS_E2E:-"false"}
ENABLE_GKE_E2E=${ENABLE_GKE_E2E:-"false"}
ENABLE_AKS_E2E=${ENABLE_AKS_E2E:-"false"}
ENABLE_UNIT_TESTS=${ENABLE_UNIT_TESTS:-true}
PATH=$PATH:${GOPATH}/bin:${PWD}/bin

mkdir -p ${GOPATH}/bin

usage() {
  cat <<EOF
  Usage: $(basename $0)
  --enable-conformance   : run the kubernetes conformance check suite
  --enable-unit-tests    : run the bats unit tests
  --enable-e2e-user      : indicates we should provision a test user
  --enable-gke           : indicates we should run E2E on GKE (default: ${ENABLE_GKE_E2E})
  --enable-eks           : indicates we should run E2E on EKS (default: ${ENABLE_EKS_E2E})
  --enable-aks           : indicates we should run E2E on AKS (default: ${ENABLE_AKS_E2E})
  --e2e-team             : is the name of the team to use
  -h|--help              : display this usage menu
EOF
  if [[ -n $@ ]]; then
    echo "[error] $@"
    exit 1
  fi
  exit 0
}

create-aws-profile() {
  aws configure --profile "${AWS_KORE_PROFILE}" set aws_access_key_id "${AWS_ACCESS_KEY_ID}"
  aws configure --profile "${AWS_KORE_PROFILE}" set aws_secret_access_key "${AWS_SECRET_ACCESS_KEY}"
  aws configure --profile "${AWS_KORE_PROFILE}" set default.region "${AWS_DEFAULT_REGION}"
}

enable-admin-user() {
  announce "Provisioning a admin e2e user: ${ADMIN_USER}"
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
  --enable-conformance) ENABLE_CONFORMANCE="true";      shift 1; ;;
  --enable-unit-tests)  ENABLE_UNIT_TESTS=${2:-"true"}; shift 1; ;;
  --enable-gke)         ENABLE_GKE_E2E=${2};            shift 2; ;;
  --enable-eks)         ENABLE_EKS_E2E=${2};            shift 2; ;;
  --enable-aks)         ENABLE_AKS_E2E=${2};            shift 2; ;;
  --enable-e2e-user)    ADMIN_USER=$2;                  shift 2; ;;
  --e2e-team)           TEAM=$2;                        shift 2; ;;
  -h|--help)            usage;                                   ;;
  *)                                                    shift 1; ;;
  esac
done

if [[ -n "${ADMIN_USER}" ]]; then
  enable-admin-user
fi

if [[ "${ENABLE_UNIT_TESTS}" == "true" ]]; then
  mkdir -p ${E2E_DIR}
  # @step: write the credentials to disk
  if [[ -n "${GKE_SA_QA}" ]]; then
    echo -n ${GKE_SA_QA} | base64 -d > ${E2E_DIR}/gke-credentials.yml 2>/dev/null
  fi
  if [[ -n "${EKS_SA_QA}" ]]; then
    echo -n ${EKS_SA_QA} | base64 -d > ${E2E_DIR}/eks-credentials.yml 2>/dev/null
  fi
  if [[ -n "${AKS_SA_QA}" ]]; then
    echo -n ${AKS_SA_QA} | base64 -d > ${E2E_DIR}/aks-credentials.yml 2>/dev/null
  fi

  if [[ "${ENABLE_GKE_E2E}" == "true" ]]; then
    if [[ ! -f "${E2E_DIR}/gke-credentials.yml" ]]; then
      error "you need to have the ${E2E_DIR}/gke-credentials.yml containing real credentials for gke"
      error "you can copy an template from: examples/gcp-credentials.yml"
      exit 1
    fi
  fi

  if [[ "${ENABLE_AKS_E2E}" == "true" ]]; then
    if [[ ! -f "${E2E_DIR}/aks-credentials.yml" ]]; then
      error "you need to have the ${E2E_DIR}/aks-credentials.yml containing real credentials for aks"
      error "you can copy an template from: examples/aks-credentials.yml"
      exit 1
    fi
  fi

  if ! test/e2e/create-kore-config.sh; then
    error "failed to generate a client configuration for cli"
    exit 1
  fi

  if ! create-aws-profile; then
    error "failed to create the aws profile"
    exit 1
  fi

  if ! test/e2e/check-suite-units.sh \
    --enable-gke ${ENABLE_GKE_E2E} \
    --enable-eks ${ENABLE_EKS_E2E} \
    --enable-aks ${ENABLE_AKS_E2E}; then
    error "failed trying on unit tests"
    exit 1
   fi
else
  announce "skipping the unit tests suite"
fi

if [[ "${ENABLE_CONFORMANCE}" == "true" ]]; then
  test/e2e/check-suite-conformance.sh || exit 1
else
  announce "skipping the kubernetes e2e conformance suite"
fi
