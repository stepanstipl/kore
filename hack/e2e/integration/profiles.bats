#!/bin/bash
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

@test "We should be able to list our profiles" {
  run bash -c "${KORE} profiles ls | grep local"
  [[ "$status" -eq 0 ]]
  run bash -c "${KORE} profiles list | grep local"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to show the current profile" {
  run bash -c "${KORE} profiles show | grep 127.0.0.1:10080"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to use a profile" {
  run bash -c "${KORE} profiles use local"
  [[ "$status" -eq 0 ]]
  run bash -c "${KORE} profiles show | grep 127.0.0.1:10080"
  [[ "$status" -eq 0 ]]
}

@test "We should fail if the profile does not exist" {
  run bash -c "${KORE} profiles use not_there || true"
  [[ "$status" -eq 0 ]]
}

@test "We should fail trying to delete a profile that doesn't exist" {
  run bash -c "${KORE} profiles delete not_there || true"
  [[ "$status" -eq 0 ]]
}

