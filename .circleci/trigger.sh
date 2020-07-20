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

: ${CIRCLECI_TOKEN?"Your circle ci token must be set in the environment to trigger jobs"}

BRANCH=$(git rev-parse --abbrev-ref HEAD)

usage() {
  cat <<EOF
  Usage: $(basename $0)
  --branch <name>   : the branch to run the build on (defaults: ${BRANCH})
  -h|--help         : display this usage menu
EOF
  if [[ -n $@ ]]; then
    echo "[error] $@"
    exit 1
  fi
  exit 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
  --branch)     BRANCH=${2};         shift 2; ;;
  -h|--help)    usage;                        ;;
  *)                                 shift 1; ;;
  esac
done

echo "Attempting to trigger the build, branch: ${BRANCH}"

JOB=$(curl -s -u ${CIRCLECI_TOKEN}: -X POST \
  --header "Accept: application/json" \
  --header "Content-Type: application/json" -d "{
    \"branch\": \"${BRANCH}\",
    \"parameters\": {
    }
  }" \
  https://circleci.com/api/v2/project/github/appvia/kore/pipeline)

if [[ $? -ne 0 ]]; then
  echo "[error] failed to trigger the pipeline, ${JOB}"
  exit 1
fi

if grep -q "pending" <<<"${JOB}"; then
  echo "Pipeline appears to be been triggered"
else
  echo "Pipeline doesn't appear to have been triggered: ${JOB}"
  exit 1
fi
