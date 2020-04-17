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

@test "We should able to delete the ${TEAM} team now" {
  runit "${KORE} delete team ${TEAM}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get teams | grep ^${TEAM} || true"
  [[ "$status" -eq 0 ]]
}

@test "We should able to delete the ${TEAM} gke credentials" {
  runit "${KORE} delete -f ${BASE_DIR}/${E2E_DIR}/gke-credentials.yml -t kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to delete the user from kore" {
  runit "${KORE} delete user test"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get users | grep ^test || true"
  [[ "$status" -eq 0 ]]
}
