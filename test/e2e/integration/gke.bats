#!/usr/bin/env bats
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
load helper

@test "Ensuring we have an allocation to build a cluster in GKE" {
  runit "${KORE} get allocations -t ${TEAM} | grep ^gke"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to build a cluster ${CLUSTER} in GKE" {
  if runit "${KORE} get clusters ${CLUSTER} -t ${TEAM}"; then
    skip
  else
    runit "${KORE} create cluster -p gke-development -a gke ${CLUSTER} --show-time -t e2e"
    [[ "$status" -eq 0 ]]
  fi
}

@test "We should be able to see the cluster ${CLUSTER} in the team" {
  runit "${KORE} get clusters ${CLUSTER} -t ${TEAM} | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "The user should be able to see the endpoint" {
  runit "${KORE} get clusters ${CLUSTER} -t ${TEAM} -o json | jq '.status.endpoint' | grep null || true"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to see the gkes cloud provider in the team" {
  runit "${KORE} get gkes ${CLUSTER} -t ${TEAM}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get gkes ${CLUSTER} -t ${TEAM} | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "The cluster should have a secret related to the cluster in the team" {
  runit "${KORE} get secrets ${CLUSTER} -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}

@test "The cluster secret should contain a number of fields (token endpoint ca.crt)"  {
  for key in token endpoint ca.crt; do
    runit "${KORE} get secrets ${CLUSTER} -t ${TEAM} -o json | jq \".spec.data.${key}\" | grep null || true"
    [[ "$status" -eq 0 ]]
  done
}

@test "We should be able to generate the kubeconfig for the cluster" {
  runit "${KORE} kubeconfig -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}

@test "We should have a kubeconfig file in the home directory" {
  [[ -f "${HOME}/.kube/config" ]] && true
  [[ "$status" -eq 0 ]]
}

@test "The koreconfig file should contain a reference to the cluster" {
  grep ${CLUSTER} ${HOME}/.kube/config
  [[ "$status" -eq 0 ]]
}

@test "You should be able to retrieve the nodes of the cluster" {
  runit "${KUBECTL} --context=${CLUSTER} get nodes"
  [[ "$status" -eq 0 ]]
}

@test "We need to ensure the default psp has been created" {
  runit "${KUBECTL} --context=${CLUSTER} get psp kore.default"
  [[ "$status" -eq 0 ]]
}

@test "We should have a namespace called kore" {
  runit "${KUBECTL} --context=${CLUSTER} get namespace kore"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to run a pod on the cluster" {
  if ${KUBECTL} --context=${CLUSTER} get deployment web; then
    skip
  fi
  ${KUBECTL} --context=${CLUSTER} create deployment web --image=nginx
  [[ "$status" -eq 0 ]]
}
