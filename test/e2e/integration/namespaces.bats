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

setup() {
  ${KORE} get cluster ${CLUSTER} | grep -i deleting && skip || true
  ${KORE} get cluster ${CLUSTER} | grep -i pending && skip || true
}

@test "We should fail to provision a namespace if the cluster does not exist" {
  run bash -c "${KORE} create namespace nothing -c does_not_exist -t ${TEAM}"
  [[ "$status" -eq 1 ]]
}

@test "We want to ensure that the namespaces are provisioned" {
  namespace="ingress"
  fullname="${CLUSTER}-${namespace}"

  if ! ${KORE} get namespaceclaims ${fullname} -t ${TEAM} 2>/dev/null; then
    runit "${KORE} create namespace -c ${CLUSTER} ${namespace} -t ${TEAM}"
    [[ "$status" -eq 0 ]]
  fi

  retry 4 "${KORE} get namespaceclaims ${fullname} -t ${TEAM} -o json | jq '.status.status' | grep -i success"
  [[ "$status" -eq 0 ]]
}

@test "We should have a namespace which is tagged by owned" {
  namespace="ingress"

  runit "${KUBECTL} --context=${CLUSTER} get namespace ${namespace} -o json | jq -r '.metadata.label[\"kore.appvia.io/owned\"]' | grep null || true"
  [[ "$status" -eq 0 ]]
}

@test "We need to ensure the rbac rolebinding kore:team has been provision in the namespace" {
  name="ingress"

  runit "${KUBECTL} --context=${CLUSTER} -n ${name} get rolebinding kore:team"
  [[ "$status" -eq 0 ]]
}

@test "Ensure the role binding is mapped to kore-nsadmin" {
  name="ingress"

  runit "${KUBECTL} --context=${CLUSTER} -n ${name} get rolebinding kore:team -o json | jq '.roleRef.name' | grep kore-nsadmin"
  [[ "$status" -eq 0 ]]
}

@test "Ensure the role binding has the team member" {
  name="ingress"
  username=$(${KORE} whoami | tail -n1 | awk '{ print $1 }')
  [[ "$status" -eq 0 ]]

  runit "${KUBECTL} --context=${CLUSTER} -n ${name} get rolebinding kore:team -o yaml | grep 'name: ${username}'"
  [[ "$status" -eq 0 ]]
}

@test "Ensure when adding a member to the team the rbac is updated" {
  runit "${KORE} create member -u test -t ${TEAM}"
  [[ "$status" -eq 0 ]]
  runit "${KORE} get member -t ${TEAM} | grep ^test"
  [[ "$status" -eq 0 ]]

  name="ingress"

  retry 4 "${KUBECTL} --context=${CLUSTER} -n ${name} get rolebinding kore:team -o yaml | grep 'name: test'"
  [[ "$status" -eq 0 ]]
}

@test "Ensure when the user is deleted from the group the rbac rule is ammended" {
  runit "${KORE} delete member -u test -t ${TEAM}"
  [[ "$status" -eq 0 ]]

  name="ingress"

  retry 4 "${KUBECTL} --context=${CLUSTER} -n ${name} get rolebinding kore:team -o yaml | grep 'name: test' || true"
  [[ "$status" -eq 0 ]]
}

@test "We should be able to delete the namespace from the cluster" {
  namespace="ingress"
  fullname="${CLUSTER}-${namespace}"

  runit "${KORE} delete namespaceclaims ${fullname} -t ${TEAM}"
  [[ "$status" -eq 0 ]]
  retry 10 "${KORE} get namespaceclaims ${fullname} -t ${TEAM} 2>&1 | grep 'not found$'"
  [[ "$status" -eq 0 ]]
  retry 10 "${KUBECTL} --context=${CLUSTER} get namespace ${namespace} 2>&1 | grep -q 'not found$'"
  [[ "$status" -eq 0 ]]
}
