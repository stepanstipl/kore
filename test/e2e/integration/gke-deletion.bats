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

@test "We should be able to delete the gke cluster" {
  ${KORE} get cluster ${CLUSTER} -t ${TEAM} || skip

  retry 10 "${KORE} delete cluster ${CLUSTER} -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}

@test "The status of the cluster should change to deleting" {
  ${KORE} get cluster ${CLUSTER} -t ${TEAM} || skip

  retry 10 "${KORE} get cluster ${CLUSTER} -t ${TEAM} -o json | jq '.status.status' | grep -i deleting"
  [[ "$status" -eq 0 ]]
}

@test "We should see the status of the gke cluster change to deleting" {
  ${KORE} get cluster ${CLUSTER} -t ${TEAM} || skip

  retry 10 "${KORE} get gkes ${CLUSTER} -t ${TEAM} -o json | jq '.status.status' | grep -i deleting"
  [[ "$status" -eq 0 ]]
}

@test "We should see the gkes resource be deleted" {
  ${KORE} get cluster ${CLUSTER} -t ${TEAM} || skip

  retry 120 "${KORE} get gkes ${CLUSTER} -t ${TEAM} 2>&1 | grep 'does not exist'"
  [[ "$status" -eq 0 ]]
}

@test "We should see the cluster resource be deleted" {
  ${KORE} get cluster ${CLUSTER} -t ${TEAM} || skip

  retry 120 "${KORE} get clusters ${CLUSTER} -t ${TEAM} 2>&1 | grep 'does not exist'"
  [[ "$status" -eq 0 ]]
}

