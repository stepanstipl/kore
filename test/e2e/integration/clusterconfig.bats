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

@test "We should find the cluster config in the kore namespace" {
  runit "${KUBECTL} --context=${CLUSTER} get namespace kore"
  [[ "$status" -eq 0 ]]
  runit "${KUBECTL} --context=${CLUSTER} -n kore get secret kore-config"
  [[ "$status" -eq 0 ]]
}

@test "We should find the client certificate in the cluster config secret" {
  for key in api-url hub-url idp-server-url tls.crt tls.key; do
    runit "${KUBECTL} --context=${CLUSTER} -n kore get secret -o json | jq -r \".data.${key}\" | grep null || true"
    [[ "$status" -eq 0 ]]
  done
}

