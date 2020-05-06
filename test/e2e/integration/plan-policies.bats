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

@test "We should find polcies for eks and gke" {
  runit "${KORE} get planpolicies"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get planpolicies gke-default"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get planpolicies eks-default"
  [[ "$status" -eq 0 ]]
}

@test "We should see the default policies are allocated to all teams" {
  runit "${KORE} get allocations planpolicy-default-eks -t ${TEAM}"  
  [[ "$status" -eq 0 ]]
  runit "${KORE} get allocations planpolicy-default-gke -t ${TEAM}"  
  [[ "$status" -eq 0 ]]
}