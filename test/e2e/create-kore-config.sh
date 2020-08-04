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

source test/e2e/environment.sh || exit 1

KORE_CONFIG=${HOME}/.kore/config
KORE_API_URL=${KORE_API_PUBLIC_URL_E2E:-"http://localhost:10080"}
KORE_IDP_SERVER_URL=${KORE_IDP_SERVER_URL:-"unknown"}
KORE_IDP_CLIENT_ID=${KORE_IDP_CLIENT_ID:-"unknown"}
KORE_ID_TOKEN=${KORE_ID_TOKEN_QA:-"unknown"}

[[ -f ${KORE_CONFIG} ]] && exit 0

mkdir -p $(dirname ${KORE_CONFIG}) || {
  error "unable to create the client configuration directory";
  exit 1;
}

announce "Generating a kore configuration: ${KORE_CONFIG}"
cat << EOF > ${KORE_CONFIG}
current-profile: local
profiles:
  local:
    server: local
    user: local
servers:
  local:
    server: ${KORE_API_URL}
users:
  local:
    oidc:
      authorize-url: ${KORE_IDP_SERVER_URL}
      client-id: ${KORE_IDP_CLIENT_ID}
      id-token: ${KORE_ID_TOKEN}
EOF