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

setup() {
  ${KORE} get cluster ${CLUSTER} | grep -i deleting && skip || true
  ${KORE} get cluster ${CLUSTER} | grep -i pending && skip || true
}

@test "We should be able to retrieve the namespaces from the cluster" {
  runit "${KUBECTL} --context=${CLUSTER} get ns"
  [[ "$status" -eq 0 ]]
}

@test "We should not be about to access to the cluster with an invalid token" {
  runit "${KUBECTL} --context=${CLUSTER} --token=invalid get no 2>&1 | grep '^Error from server (Forbidden)'"
  [[ "$status" -eq 0 ]]
}

@test "We should be about to access the cluster via a valid kubernetes token" {
  SA="kube-test"
  if ! ${KUBECTL} --context=${CLUSTER} get sa ${SA}; then
    runit "${KUBECTL} --context=${CLUSTER} create sa ${SA}"
    [[ "$status" -eq 0 ]]
    runit "${KUBECTL} --context=${CLUSTER} create rolebinding --clusterrole=view --serviceaccount=default:${SA} ${SA}"
    [[ "$status" -eq 0 ]]
  fi
  runit "${KUBECTL} --context=${CLUSTER} get sa ${SA} -o json | jq -r '.secrets[0].name' > /tmp/default.sa"
  [[ "$status" -eq 0 ]]
  runit "${KUBECTL} --context=${CLUSTER} get secret $(cat /tmp/default.sa) | jq -r '.data.token' | base64 -d > /tmp/default.token"
  [[ "$status" -eq 0 ]]
  runit "${KUBECTL} --context=${CLUSTER} --token=$(cat /tmp/default.token) get po"
  [[ "$status" -eq 0 ]]
  runit "${KUBECTL} --context=${CLUSTER} --token=$(cat /tmp/default.token) get node || false"
  [[ "$status" -eq 0 ]]
}

@test "If we change the auth-proxy allowed range we should lose access to the cluster" {
  if ! ${KORE} get clusters ${CLUSTER} -t ${TEAM} -o yaml | grep 1.1.1.1; then
    runit "${KUBECTL} --context=${CLUSTER} get nodes"
    [[ "$status" -eq 0 ]]
    runit "${KORE} alpha patch clusters ${CLUSTER} spec.configuration.authProxyAllowedIPs.0 1.1.1.1/32 -t ${TEAM}"
    [[ "$status" -eq 0 ]]
  fi
  retry 10 "${KUBECTL} --context=${CLUSTER} get nodes 2>&1 | grep '^Error from server (Forbidden)'"
  [[ "$status" -eq 0 ]]
}

@test "If we revert the allowed network range back, we should see the cluster again" {
  runit "${KORE} alpha patch clusters ${CLUSTER} spec.configuration.authProxyAllowedIPs.0 0.0.0.0/0 -t ${TEAM}"
  [[ "$status" -eq 0 ]]
  retry 10 "${KUBECTL} --context=${CLUSTER} get nodes 2>&1"
  [[ "$status" -eq 0 ]]
}
