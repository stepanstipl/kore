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
  if ${KORE} get cluster ${CLUSTER} -o json | jq '.status.status' | grep -i deleting; then
    skip "Cluster is already deleting, skipping these checks"
  fi
}

@test "We should be able to apply the EKS credentials" {
  runit "${KORE} apply -f ${BASE_DIR}/e2eci/eks-credentials.yml -t kore-admin"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get ekscredentials aws -t kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "We should have an allocation for EKS credentials" {
  runit "${KORE} get allocations aws -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to build a cluster in EKS" {
  runit "${KORE} create cluster ${CLUSTER} -t ${TEAM} -a aws -p eks-development --no-wait"
  [[ "$status" -eq 0 ]]
}

@test "We should see a VPC provisioned for us" {
  retry 100 "${KORE} get eksvpc ${CLUSTER} -t ${TEAM} -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should see a EKS cluster provisions" {
  retry 300 "${KORE} get eks ${CLUSTER} -t ${TEAM} -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should have a default EKS Nodegroup provision" {
  retry 240 "${KORE} get eksnodegroup ${CLUSTER}-default -t ${TEAM} -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should see the cluster go successful" {
  retry 60 "${KORE} get cluster ${CLUSTER} -t ${TEAM} -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "The cluster secret should contain a number of fields (token endpoint ca.crt)"  {
  runit "${KORE} get secrets ${CLUSTER} -t ${TEAM}"
  [[ "$status" -eq 0 ]]
  for key in token endpoint ca.crt; do
    runit "${KORE} get secrets ${CLUSTER} -t ${TEAM} -o json | jq \".spec.data.${key}\" | grep null || true"
    [[ "$status" -eq 0 ]]
  done
}

@test "We should be able to generate the kubeconfig for the cluster" {
  runit "${KORE} kubeconfig -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}

@test "You should be able to retrieve the nodes of the cluster" {
  retry 60 "${KUBECTL} --context=${CLUSTER} get nodes"
  [[ "$status" -eq 0 ]]
}

@test "We should have a namespace called kore" {
  runit "${KUBECTL} --context=${CLUSTER} get namespace kore"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to run a pod on the cluster" {
  if ! ${KUBECTL} --context=${CLUSTER} get deployment web; then
    runit "${KUBECTL} --context=${CLUSTER} create deployment web --image=nginx"
    [[ "$status" -eq 0 ]]
  fi

  runit "${KUBECTL} --context=${CLUSTER} get deployment web"
  [[ "$status" -eq 0 ]]
  runit "${KUBECTL} --context=${CLUSTER} get deployment web"
  [[ "$status" -eq 0 ]]
  runit "${KUBECTL} --context=${CLUSTER} get pod | grep '^web.*Running'"
  [[ "$status" -eq 0 ]]
}

@test "We should have a default pod security policy eks.privileged" {
  runit "${KUBECTL} --context=${CLUSTER} get psp eks.privileged"
  [[ "$status" -eq 0 ]]
}

@test "We should see the number of nodes change when I update the desired state" {
  runit "${KORE} alpha patch cluster ${CLUSTER} spec.nodeGroups.0.desiredSize 2"
  [[ "$status" -eq 0 ]]
  retry 60 "${KUBECTL} --context=${CLUSTER} get nodes --no-headers | grep Ready | wc -l | grep 2"
  [[ "$status" -eq 0 ]]
}

