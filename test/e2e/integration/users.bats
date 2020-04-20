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

@test "Ensuring the admin user is showing up" {
  retry 3 "${KORE} get user admin"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring you cannot delete the admin user" {
  runit "${KORE} delete user admin || true"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring I show up in the admin group of whoami" {
  runit "${KORE} whoami | grep kore-admin"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring I show I don't show up in kore-default" {
  runit "${KORE} whoami | grep kore-default || true"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring we can create a new users from examples" {
  runit "${KORE} apply -f ../../../examples/user.yml"
  [[ "$status" -eq 0 ]]
}

@test "Ensure the user shows up in the list" {
  runit "${KORE} get user test"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring the new user is not disabled" {
  runit "${KORE} get user test | awk '/test/ { print $3 }' | grep false"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to add the user to the kore-default team" {
  runit "${KORE} create member -u test -t kore-default"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get member -t kore-default | grep test"
  [[ "$status" -eq 0 ]]
}

@test "Ensuring the user can be added to the team multiple times without error" {
  runit "${KORE} create member -u test -t kore-default"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to delete member from the kore-default " {
  runit "${KORE} delete member -u test -t kore-default"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get member -t kore-default | grep ^test || true"
  [[ "$status" -eq 0 ]]
}

@test "We should not be allowed to remove the admin user from the kore-admin team" {
  runit "${KORE} delete member -u admin -t kore-admin || true"
  [[ "$status" -eq 0 ]]
}

@test "We should to create a user called e2e" {
  runit "${KORE} create user --email e2e@appvia.io e2e@appvia.io"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to add the user to the admin group" {
  runit "${KORE} create admin -u e2e@appvia.io"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get admin | grep ^e2e@appvia.io"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to delete the user from admin group" {
  runit "${KORE} delete admin -u e2e@appvia.io"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get admin | grep ^e2e@appvia.io || true"
  [[ "$status" -eq 0 ]]
  runit "${KORE} delete user e2e@appvia.io"
  [[ "$status" -eq 0 ]]
}
