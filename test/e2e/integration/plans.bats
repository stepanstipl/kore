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

@test "We should be able to list the embedded plans" {
  runit "${KORE} get plans"
  [[ "$status" -eq 0 ]]
}

@test "We should have a gke-development plan" {
  runit "${KORE} get plans gke-development"
  [[ "$status" -eq 0 ]]
}

@test "We should have a gke production plan" {
  runit "${KORE} get plans gke-production"
  [[ "$status" -eq 0 ]]
}

@test "The plans should include valid json data" {
  runit "${KORE} get plans gke-development -o json | jq"
  [[ "$status" -eq 0 ]]
}

@test "We should see a valid version in the gke plan" {
  runit "${KORE} get plans gke-development -o json | jq '.spec.values.version' | grep gke"
  [[ "$status" -eq 0 ]]
}

