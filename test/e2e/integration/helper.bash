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

# Retry a command $1 times until it succeeds. Wait $2 seconds between retries.
retry() {
  local attempts=$1; shift;
  local delay=5
  local i

  for ((i=0; i < attempts; i++)); do
      #echo "Executing: $@" >> command.log
      run bash -c "$@"
      if [[ "${status}" -eq 0 ]]; then
        echo "$output"
        return 0
      fi
      sleep $delay
  done

  echo "Command \"$@\" failed $attempts times. Status: $status. Output: $output" >&2
  false
}

runit() {
  retry 5 "$@"
}

# wait-on-deployment is responsible for waiting for a deployment to deploy
wait-on-deployment() {
  local namespace=$1
  local labels=$2
  local expected=${3:-1}

  for ((i=0; i<60; i++)) do
    count=$(kubectl -n ${namespace} get po -l ${labels} --field-selector=status.phase=Running --no-headers | grep -i running | wc -l)
    if [[ $? -eq 0 ]]; then
      if [[ "${count}" == "${expected}" ]]; then
        return 0
      fi
    fi
    sleep 2
  done

  false
}
