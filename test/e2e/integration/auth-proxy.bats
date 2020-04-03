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

@test "If we change the auth-proxy allowed range we should lose access to the cluster" {
  tempfile="${BASE_DIR}/${E2E_DIR}/gke.auth"

  if ! ${KORE} get clusters ${CLUSTER} -t e2e -o yaml | grep 1.1.1.1; then
    runit "${KUBECTL} --context=${CLUSTER} get nodes"
    [[ "$status" -eq 0 ]]
    runit "${KORE} get clusters ${CLUSTER} -t e2e -o yaml > ${tempfile}"
    [[ "$status" -eq 0 ]]
    runit "sed -i -e 's/\[\"0.0.0.0\/0\"\]/[\"1.1.1.1\/32\"]/' ${tempfile}"
    [[ "$status" -eq 0 ]]
    runit "${KORE} apply -f ${tempfile} -t e2e"
    [[ "$status" -eq 0 ]]
  fi
  retry 10 "${KUBECTL} --context=${CLUSTER} get nodes || true"
  [[ "$status" -eq 0 ]]
}

@test "If we revert the allowed network range back, we should see the cluster again" {
  tempfile=${BASE_DIR}/${E2E_DIR}/gke.auth

  runit "${KORE} get clusters ${CLUSTER} -t e2e -o yaml > ${tempfile}"
  [[ "$status" -eq 0 ]]
  runit "sed -i -e 's/\[\"1.1.1.1\/32\"\]/[\"0.0.0.0\/0\"]/' ${tempfile}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} apply -f ${tempfile} -t e2e"
  [[ "$status" -eq 0 ]]
  retry 20 "${KUBECTL} --context=${CLUSTER} get nodes"
  [[ "$status" -eq 0 ]]
  runit "rm -f ${tempfile} || false"
  [[ "$status" -eq 0 ]]
}
