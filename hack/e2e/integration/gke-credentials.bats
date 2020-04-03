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

@test "We should not have any GKE credentials existing" {
  runit "${KORE} delete gkecredentials gke -t kore-admin || true"
  [[ "$status" -eq 0 ]]
}

@test "We should not be able to list any gke credentials" {
  runit "${KORE} get gkecredentials gke -t kore-admin || true"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to apply the gke credentials" {
  run ${KORE} apply -f ${BASE_DIR}/e2eci/gke-credentials.yml -t kore-admin
  [[ "$status" -eq 0 ]]
}

@test "We should be able to view the credentials via the cli now" {
  run ${KORE} get gkecredentials gke -t kore-admin
  [[ "$status" -eq 0 ]]
}

@test "The GKE credentials should come back as verified" {
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.verified' | grep -i true"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "If we break the GKE credentials the verification should fail" {
  run ${KORE} apply -f ${BASE_DIR}/examples/gcp-credentials.yml -t kore-admin
  [[ "$status" -eq 0 ]]
}

@test "The GKE credentials should come back as failed" {
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.verified' | grep -i false"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.status' | grep -i failure"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to reapply the valid credentials and be ok" {
  run ${KORE} apply -f ${BASE_DIR}/e2eci/gke-credentials.yml -t kore-admin
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.verified' | grep -i true"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to see the gke allocations in the e2e team" {
  runit "${KORE} get allocations gke -t e2e"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get allocations gke -o json -t e2e | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "If we delete the allocation, the e2e should no longer see the gke credentials" {
  runit "${KORE} get allocations gke -t kore-admin"
  [[ "$status" -eq 0 ]]
  runit "${KORE} delete allocations gke -t kore-admin"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get allocations -t e2e | grep ^gke || true"
  [[ "$status" -eq 0 ]]
}

@test "We should reapply the credentials and get back the allocation in the e2e team" {
  runit "${KORE} apply -f ${BASE_DIR}/e2eci/gke-credentials.yml -t kore-admin"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.verified' | grep -i true"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get gkecredentials gke -t kore-admin -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
  retry 5 "${KORE} get allocations gke -t e2e"
  [[ "$status" -eq 0 ]]
}
