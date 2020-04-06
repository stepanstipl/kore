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

[[ "${DEBUG}" == "true" ]] && set -x

export PLATFORM_DIR="test/e2e"
export BASE_DIR="../../.."
export KUBECTL="kubectl"
export KORE="korectl"
export NC='\e[0m'
export GREEN='\e[0;32m'
export YELLOW='\e[0;33m'
export RED='\e[0;31m'

# log will echo the message and reset colour code
log()      { (2>/dev/null echo -e "$@${NC}"); }
# announce will log a given message (used for standard info logging)
announce() { log "${GREEN}[$(date)][INFO] $@"; }
# failed is to notify of configuration failures (e.g. missing files / environment variables)
failed()   { log "${YELLOW}[$(date)][FAIL] $@"; }
# error is used when unexpected errors occur (e.g. unable to communicate with API)
error()    { log "${RED}[$(date)][ERROR] $@"; }

# attempt tries to perform a command x number of times before giving up
attempt() {
  max_attempts=2
  for ((i=1; i<=${max_attempts}; i++)); do
    if eval "$@"; then
      return 0
    else
      error "failed to run command: '${@}', retrying (attempt/max = ${i}/${max_attempts})"
      sleep 5
    fi
  done
  return 1
}
