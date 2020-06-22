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

@test "We should be able to create a secret" {
  runit "${KORE} create secret test --from-literal=hello=world -d test -t ${TEAM} --force"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get secret test -o yaml -t ${TEAM} | grep hello:"
  [[ "$status" -eq 0 ]]
}

@test "We should not be able to create a secret which exists" {
  runit "${KORE} create secret test --from-literal=hello=world -d test -t ${TEAM} 2>&1 | grep 'already exists'"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to force a secert to overwrite" {
  runit "${KORE} create secret test --from-literal=hello=world -d test -t ${TEAM} --force"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to delete the secret" {
  if ! ${KORE} get secret -t ${TEAM}; then
    skip "secret does not exist"
  fi
  runit "${KORE} delete secret test -t ${TEAM}"
  [[ "$status" -eq 0 ]]
}

