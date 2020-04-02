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

@test "Ensuring the admin user is showing up" {
  retry 3 "${KORE} get user admin"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring you cannot delete the admin user" {
  retry 1 ${KORE} delete user admin
  [[ "$status" -eq 1 ]]
}

@test "Ensuring I show up in the admin group of whoami" {
  retry 2 "${KORE} whoami | grep kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring I show you do not show up in kore-default anymore" {
  retry 2 "${KORE} whoami | grep kore-default || true"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring we can create a new users from examples" {
  retry 2 "${KORE} apply -f ../../../examples/user.yml"
  [[ "$status" -eq 0 ]]
}

@test "Ensure the user shows up in the list" {
  retry 2 "${KORE} get user test"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring the new user is not disabled" {
  retry 2 "${KORE} get user test | awk '/test/ { print $3 }' | grep false"
  [[ "$status" -eq 0 ]]
}

@test "Checking we can add the user to the kore-default team" {
  retry 2 "${KORE} create member -u test -t kore-default"
  [[ "$status" -eq 0 ]]
  retry 2 "${KORE} get member -t kore-default | grep test"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring the user can be added to the team multiple times without error" {
  retry 2 "${KORE} create member -u test -t kore-default"
  [[ "$status" -eq 0 ]]
}

