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

set -o errexit
set -o nounset
set -o pipefail

: ${KORE_API_PUBLIC_URL_QA?"The QA Kore API URL must be set"}
: ${KORE_UI_PUBLIC_URL_QA?"The QA Kore UI URL must be set"}

export ENVIRONMENT="${ENVIRONMENT:-"qa"}"
export BUILD_ID=${BUILD_ID:-${VERSION}}
export KORE_UI_SHOW_PROTOTYPES="true"
export VERSION=${VERSION:-"latest"}
export KORE_API_PUBLIC_URL=${KORE_API_PUBLIC_URL_QA}
export KORE_UI_PUBLIC_URL=${KORE_UI_PUBLIC_URL_QA}

hack/deploy.sh
