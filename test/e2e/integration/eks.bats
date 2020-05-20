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
  runit "${KORE} apply -f ${BASE_DIR}/e2eci/eks-credentials.yml -t kore-admin 2>&1 >/dev/null"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get ekscredentials aws -t kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "We should have an allocation for EKS credentials" {
  runit "${KORE} get allocations ekscredentials-aws -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to build a cluster in EKS" {
  ${KORE} get cluster ${CLUSTER} -t ${TEAM} && skip

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

@test "We should find two iam policies related to the cluster" {
 runit "aws --profile ${AWS_KORE_PROFILE} iam list-roles | jq '.Roles[].RoleName' -r | grep ${CLUSTER}"
  [[ "$status" -eq 0 ]]
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
  runit "${KORE} alpha patch cluster ${CLUSTER} spec.configuration.nodeGroups.0.desiredSize 2 -t ${TEAM}"
  [[ "$status" -eq 0 ]]
  retry 60 "${KUBECTL} --context=${CLUSTER} get nodes --no-headers | grep Ready | wc -l | grep 2"
  [[ "$status" -eq 0 ]]
}

#
### Disabling for now as this add almost 15 minutes to the E2E .. honestly how can it take
# 15 MINUTES to update a firewall rule
#
#@test "We should be able to update the authorized master list of eks and lose access" {
#  cat <<EOF > /tmp/plan-policy
#---
#apiVersion: config.kore.appvia.io/v1
#kind: PlanPolicy
#metadata:
#  name: eks-open
#spec:
#  description: Allows changing the AuthorizedMasterNetworks
#  summary: Allows changing the AuthorizedMasterNetworks
#  kind: EKS
#  properties:
#  - allowUpdate: true
#    disallowUpdate: false
#    name: authorizedMasterNetworks
#  - allowUpdate: true
#    disallowUpdate: false
#    name: authProxyAllowedIPs
#  - allowUpdate: true
#    disallowUpdate: false
#    name: clusterUsers
#  - allowUpdate: true
#    disallowUpdate: false
#    name: defaultTeamRole
#  - allowUpdate: true
#    disallowUpdate: false
#    name: description
#  - allowUpdate: true
#    disallowUpdate: false
#    name: domain
#  - allowUpdate: true
#    disallowUpdate: false
#    name: nodeGroups
#  - allowUpdate: true
#    disallowUpdate: false
#    name: privateIPV4Cidr
#  - allowUpdate: true
#    disallowUpdate: false
#    name: region
#  - allowUpdate: true
#    disallowUpdate: false
#    name: version
#---
#apiVersion: config.kore.appvia.io/v1
#kind: Allocation
#metadata:
#  name: eks-open
#spec:
#  name: eks-open
#  resource:
#    group: config.kore.appvia.io
#    kind: PlanPolicy
#    name: eks-open
#    namespace: kore-admin
#    version: v1
#  summary: Allows changing the authorizedMasterNetworks
#  teams:
#    - '*'
#EOF
#  [[ "$status" -eq 0 ]]
#  runit "${KORE} apply -f /tmp/plan-policy -t kore-admin"
#  [[ "$status" -eq 0 ]]
#  runit "${KORE} alpha patch cluster ${CLUSTER} spec.configuration.authorizedMasterNetworks.0 1.1.1.1/32"
#  [[ "$status" -eq 0 ]]
#  retry 60 "aws eks describe-cluster --name ${CLUSTER} | jq -r '.cluster.resourcesVpcConfig.publicAccessCidrs[0]' | grep '1.1.1.1/32'"
#  [[ "$status" -eq 0 ]]
#}
#
#@test "We should be able to revert the change back and restore access" {
#  runit "${KORE} delete -f /tmp/plan-policy -t kore-admin"
#  [[ "$status" -eq 0 ]]
#  runit "${KORE} alpha patch cluster ${CLUSTER} spec.configuration.authorizedMasterNetworks.0 0.0.0/0"
#  [[ "$status" -eq 0 ]]
#  retry 60 "aws eks describe-cluster --name ${CLUSTER} | jq -r '.cluster.resourcesVpcConfig.publicAccessCidrs[0]' | grep '0.0.0.0/0'"
#  [[ "$status" -eq 0 ]]
#  runit "rm -f /tmp/plan-policy"
#  [[ "$status" -eq 0 ]]
#}
