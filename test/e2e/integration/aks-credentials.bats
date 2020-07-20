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

@test "We should not have any AKS credentials existing" {
  runit "${KORE} delete akscredentials aks -t kore-admin || true"
  [[ "$status" -eq 0 ]]
}

@test "We should not be able to list any aks credentials" {
  runit "${KORE} get akscredentials aks -t kore-admin || true"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to apply the aks credentials" {
  runit "${KORE} apply -f ${BASE_DIR}/e2eci/aks-credentials.yml -t kore-admin 2>&1 >/dev/null"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to view the credentials via the cli now" {
  run ${KORE} get akscredentials aks -t kore-admin
  [[ "$status" -eq 0 ]]
}

@test "The AKS credentials should come back as verified" {
  retry 5 "${KORE} get akscredentials aks -t kore-admin -o json | jq '.status.verified' | grep -i true"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get akscredentials aks -t kore-admin -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to see the aks allocations in the ${TEAM} team" {
  runit "${KORE} get allocations akscredentials-aks -t ${TEAM}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get allocations akscredentials-aks -o json -t ${TEAM} | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "If we delete the allocation, the ${TEAM} should no longer see the aks credentials" {
  runit "${KORE} get allocations akscredentials-aks -t kore-admin"
  [[ "$status" -eq 0 ]]
  runit "${KORE} delete allocations akscredentials-aks -t kore-admin"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get allocations -t ${TEAM} | grep ^akscredentials-aks || true"
  [[ "$status" -eq 0 ]]
}

@test "We should reapply the credentials and get back the allocation in the ${TEAM} team" {
  runit "${KORE} apply -f ${BASE_DIR}/e2eci/aks-credentials.yml -t kore-admin 2>&1 >/dev/null"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get akscredentials aks -t kore-admin -o json | jq '.status.verified' | grep -i true"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get akscredentials aks -t kore-admin -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get allocations akscredentials-aks -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}
