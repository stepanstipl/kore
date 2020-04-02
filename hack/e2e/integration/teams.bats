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
load helper.sh

@test "Ensure the test team does not exist" {
  retry 3 "${KORE} delete team test || true"
  [[ "$status" -eq 0 ]]
  retry 3 "${KORE} get team | grep test || true"
  [[ "$status" -eq 0 ]]
}

@test "Checking we have the kore-admin team" {
  retry 2 "${KORE} get teams kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "Checking we have the kore-default team" {
  retry 2 "${KORE} get teams kore-default"
  [[ "$status" -eq 0 ]]
}

@test "Checking we error on non-existing team" {
  run ${KORE} get teams not-there
  [[ "$status" -eq 1 ]]
}

@test "We should be able to create a e2e team to run test" {
  if retry 2 "${KORE} get team | grep ^e2e"; then
    skip
  fi
  retry 2 "${KORE} create team e2e"
  [[ "$status" -eq 0 ]]
}

@test "We should see the status on the e2e team is successful " {
  retry 2 "${KORE} get teams e2e -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should see a namespace created in kore for new team" {
  retry 2 "${KUBECTL} get namespace e2e >/dev/null"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to create and delete a team and see all resource dissappeared" {
  retry 2 "${KORE} create team test"
  [[ "$status" -eq 0 ]]
  retry 2 "${KORE} get team test"
  [[ "$status" -eq 0 ]]
  retry 2 "${KUBECTL} get namespace test"
  [[ "$status" -eq 0 ]]
  retry 2 "${KORE} delete team test"
  [[ "$status" -eq 0 ]]
  retry 2 "${KORE} get team 2>&1 | grep ^test || true"
  [[ "$status" -eq 0 ]]
}

@test "We should see the namespace for the test team in kore has disappeared" {
  retry 2 "${KUBECTL} get namespace test 2>&1 | egrep -i '(not found|the server doesn|terminating)'"
  [[ "$status" -eq 0 ]]
}
