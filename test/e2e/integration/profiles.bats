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

@test "We should be able to list our profiles" {
  runit "${KORE} profiles ls | grep ${KORE_PROFILE}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} profiles list | grep ${KORE_PROFILE}"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to show the current profile" {
  runit "${KORE} profiles show"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to use a profile" {
  runit "${KORE} profiles use ${KORE_PROFILE}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} profiles show | grep ${KORE_PROFILE}"
  [[ "$status" -eq 0 ]]
}

@test "We should fail if the profile does not exist" {
  runit "${KORE} profiles use not_there || true"
  [[ "$status" -eq 0 ]]
}

@test "We should fail trying to delete a profile that doesn't exist" {
  runit "${KORE} profiles delete not_there || true"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to switch to the ${KORE_PROFILE} profile" {
  runit "${KORE} profiles use ${KORE_PROFILE}"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to set the default team of the profile" {
  runit "${KORE} profiles set current.team ${TEAM}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} profiles list | awk \"/^${KORE_PROFILE}/ { print \$3 }\" | grep ^${KORE_PROFILE}"
  [[ "$status" -eq 0 ]]
}

