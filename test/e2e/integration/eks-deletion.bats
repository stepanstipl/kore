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
  ${KORE} get cluster ${CLUSTER} -t ${TEAM} | grep EKS || skip && true
}

@test "We should be able to apply the EKS credentials" {
  runit "${KORE} apply -f ${BASE_DIR}/e2eci/eks-credentials.yml -t kore-admin"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get ekscredentials aws -t kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to delete the EKS cluster" {
  runit "${KORE} delete cluster ${CLUSTER} -t ${TEAM} --no-wait"
  [[ "$status" -eq 0 ]]
  retry 30 "${KORE} get cluster ${CLUSTER} -t ${TEAM} -o json --no-wait | jq -r '.status.status' | grep -i deleting"
  [[ "$status" -eq 0 ]]
}

@test "We should see the kubernetes resource delete" {
  ${KORE} get kubernetes ${CLUSTER} -t ${TEAM} || skip

  runit "${KORE} get kubernetes ${CLUSTER} -t ${TEAM} -o json | jq -r '.status.status' | grep -i deleting"
  [[ "$status" -eq 0 ]]
  retry 30 "${KORE} get kubernetes ${CLUSTER} -t ${TEAM} 2>&1 | grep 'not found$'"
  [[ "$status" -eq 0 ]]
}

@test "We should see the eksnodegroup resource delete" {
  ${KORE} get eksnodegroup ${CLUSTER}-default -t ${TEAM} || skip

  retry 30 "${KORE} get eksnodegroup ${CLUSTER}-default -t ${TEAM} -o json | jq -r '.status.status' | grep -i deleting"
  [[ "$status" -eq 0 ]]
  retry 180 "${KORE} get eksnodegroup ${CLUSTER}-default -t ${TEAM} 2>&1 | grep 'not exist$'"
  [[ "$status" -eq 0 ]]
}

@test "We should see the eks cluster deleted" {
  retry 300 "${KORE} get cluster ${CLUSTER} -t ${TEAM} 2>&1 | grep 'not exist$'"
  [[ "$status" -eq 0 ]]
}

@test "We should able to delete the ${TEAM} eks credentials" {
  runit "${KORE} delete -f ${BASE_DIR}/${E2E_DIR}/eks-credentials.yml -t kore-admin"
  [[ "$status" -eq 0 ]]
}
