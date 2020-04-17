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

@test "Ensure the test team does not exist" {
  runit "${KORE} delete team test || true"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get team | grep ^test || true"
  [[ "$status" -eq 0 ]]
}

@test "Checking we have the kore-admin team" {
  runit "${KORE} get teams kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "Checking we have the kore-default team" {
  runit "${KORE} get teams kore-default"
  [[ "$status" -eq 0 ]]
}

@test "Checking we error on non-existing team" {
  run ${KORE} get teams not-there
  [[ "$status" -eq 1 ]]
}

@test "We should be able to create a ${TEAM} team to run test" {
  if runit "${KORE} get team | grep ^${TEAM}"; then
    skip
  fi
  runit "${KORE} create team ${TEAM}"
  [[ "$status" -eq 0 ]]
}

@test "We should not be able to delete the kore-admin team" {
  runit  "${KORE} delete team kore-admin || true"
  [[ "$status" -eq 0 ]]
}

@test "We should see the status on the ${TEAM} team is successful " {
  runit "${KORE} get teams ${TEAM} -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to create and delete a team and see all resource disappeared" {
  runit "${KORE} create team test"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get team test"
  [[ "$status" -eq 0 ]]
  runit "${KORE} delete team test"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get team 2>&1 | grep ^test || true"
  [[ "$status" -eq 0 ]]
}
