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

BATS_OPTIONS="${BATS_OPTIONS:-""}"
ENABLE_EKS_E2E=${ENABLE_EKS_E2E:-"false"}
ENABLE_GKE_E2E=${ENABLE_GKE_E2E:-"false"}
ENABLE_AKS_E2E=${ENABLE_AKS_E2E:-"false"}
CLUSTER="${CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"}"

usage() {
  cat <<EOF
  Usage: $(basename $0)
  --enable-gke <bool>  : indicates we should run E2E on GKE (default: ${ENBALE_GKE_E2E})
  --enable-eks <bool>  : indicates we should run E2E on EKS (default: ${ENABLE_EKS_E2E})
  --enable-aks <bool>  : indicates we should run E2E on AKS (default: ${ENABLE_AKS_E2E})
  -h|--help      : display this usage menu
EOF
  if [[ -n $@ ]]; then
    error "$@"
    exit 1
  fi

  exit 0
}

announce "running the integration suite"

# run-generic-checks runs a collection fo generic cli & api checks
run-generic-checks() {
  announce "running generic unit tests"
  local units=(
      setup.bats
      profiles.bats
      users.bats
      whoami.bats
      teams.bats
      plans.bats
      secrets.bats
  )
  for unit in ${units[@]}; do
    bats ${BATS_OPTIONS} ${unit} || exit 1
  done
}

# run-cluster-checks runs a series of cluster a clister
run-cluster-checks() {
  local name=${1}
  announce "running cluster checks on cluster: ${name}"
  local units=(
    clusterroles.bats
    clusterconfig.bats
    namespaces.bats
    auth-proxy.bats
    safeguards.bats
  )
  for unit in ${units[@]}; do
    CLUSTER=${name} bats ${BATS_OPTIONS} ${unit} || exit 1
  done
}

# run-gke-check runs a collection of gke checks
run-gke-checks() {
  announce "running e2e suite on gke"
  local units=(
      gke-credentials.bats
      gke.bats
  )
  local name="${CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"}"
  # run the check on the gke
  for unit in ${units[@]}; do
    CLUSTER=${name} bats ${BATS_OPTIONS} ${unit} || exit 1
  done
  # run the cluster checks on on it
  run-cluster-checks ${name}
}

# run-eks-checks runs a collection of e2e on eks
run-eks-checks() {
  announce "running e2e suite on eks"
  local units=(
      eks.bats
  )
  local name="${CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"}"
  for unit in ${units[@]}; do
    CLUSTER=${name} bats ${BATS_OPTIONS} ${unit} || exit 1
  done
  # run the cluster checks on on it
  run-cluster-checks ${name}
}

# run-aks-checks runs a collection of e2e on aks
run-aks-checks() {
  announce "running e2e suite on aks"
  local units=(
      aks-credentials.bats
      aks.bats
  )
  local name="${CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"}"
  for unit in ${units[@]}; do
    CLUSTER=${name} bats ${BATS_OPTIONS} ${unit} || exit 1
  done
  # run the cluster checks on on it
  run-cluster-checks ${name}
}

run-teardown() {
  announce "running the teardown checks"
  local units=(
      gke-deletion.bats
      eks-deletion.bats
      aks-deletion.bats
  )
  local name="${CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"}"
  for unit in ${units[@]}; do
    CLUSTER=${name} bats ${BATS_OPTIONS} ${unit} || exit 1
  done

  CLUSTER=${CLUSTER} bats ${BATS_OPTIONS} teardown.bats || exit 1
}

# run-units is responsible for running the test framework
run-units() {
  run-generic-checks
  [[ "${ENABLE_GKE_E2E}" == "true" ]] && run-gke-checks
  [[ "${ENABLE_EKS_E2E}" == "true" ]] && run-eks-checks
  [[ "${ENABLE_AKS_E2E}" == "true" ]] && run-aks-checks
  run-teardown
}

while [[ $# -gt 0 ]]; do
  case "$1" in
  --enable-gke) ENABLE_GKE_E2E=${2}; shift 2; ;;
  --enable-eks) ENABLE_EKS_E2E=${2}; shift 2; ;;
  --enable-aks) ENABLE_AKS_E2E=${2}; shift 2; ;;
  -h|--help)    usage;                        ;;
  *)                                 shift 1; ;;
  esac
done

cd ${PLATFORM_DIR}/integration || exit 1
run-units || exit 1
