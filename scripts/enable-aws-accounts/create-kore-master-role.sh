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

log()   { (2>/dev/null echo -e "$@"); }
info()  { if [[ ! -z ${QUIET:-} ]]; then return 0; fi; log "[info]  $@"; }
error() { echo "[error] $@" 1>&2; exit 1; }

[[ ${DEBUG:-} == 'true' ]] && set -x

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
MASTER_ROLE_CF_TEMPLATE_NAME=cf-kore-master-role.json
MASTER_ROLE_CF_TEMPLATE_PATH=${SCRIPT_DIR}/${MASTER_ROLE_CF_TEMPLATE_NAME}
S3_BUCKET=${S3_BUCKET:-kore-control-tower-cf-templates}
S3_REGION=$(aws configure get region)

usage() {
cat <<EOF
Usage: $(basename $0)
  --kore-user-arn     [required] kore user to grant accees to the role e.g. "arn:aws:iam::123456789:user/kore-accounts-admin-user"
  --master-role-name  [required] name to give the master role (e.g. - "kore-accounts-role-for-custom-ou")
  --s3-bucket         bucket name to upload the stack during deploy (default: "${S3_BUCKET}")
  --s3-region         region to create the bucket in if it doesn't exist
  --dry-run           will show aws commands (not actualy run them)
  --quiet | -q        will suppress logging messages
  --help | -h         display this usage menu

EOF
  if [[ -n $@ ]]; then
    echo "[error] $@"
    exit 1
  fi
  exit 0
}

run-cmd() {
  if [[ ${DRY_RUN:-} == 'true' ]]; then
    info "dry-run:\n$@"
  else
    OUTPUT=$( $@ )
  fi
  return $?
}

describe-stack() {
  aws cloudformation describe-stacks | \
    jq -r ' .Stacks | .[] | select(.StackName=="'${1?"error missing stack name"}'") | '.${2?"missing param for describe"}''
}

print-stack-outputs() {
  stackName=${1?"error missing stack name"}
  aws cloudformation describe-stacks | \
    jq -r ' .Stacks | .[] | select(.StackName=="'${stackName}'") | .Outputs | .[] | .OutputValue'
}

create-bucket-if-required() {
  if ! aws s3 ls ${S3_BUCKET} >/dev/null 2>&1 ; then
    if [[ -z ${S3_REGION:-} ]]; then
      error "--s3-region not set so can't create bucket with correct LocationConstraint"
    fi

    info "bucket ${S3_BUCKET} not found, creating ${S3_BUCKET}..."
    if ! run-cmd aws s3api create-bucket --acl private --bucket ${S3_BUCKET} --create-bucket-configuration LocationConstraint=${S3_REGION}; then
      error "bucket ${S3_BUCKET} does not exist and can't be created"
    fi
  fi
}

deploy-stack() {
  info "deploying stack - ${KORE_MASTER_ROLE_NAME}"
  run-cmd aws cloudformation deploy \
    --stack-name ${KORE_MASTER_ROLE_NAME} \
    --template-file ${MASTER_ROLE_CF_TEMPLATE_PATH} \
    --s3-bucket ${S3_BUCKET} \
    --parameter-overrides \
      KoreUserArnParameter="${KORE_USER_ARN}" KoreMasterRoleNameParameter="${KORE_MASTER_ROLE_NAME}" \
    --capabilities CAPABILITY_NAMED_IAM \
    --no-fail-on-empty-changeset
}

wait-on-stack-complete-or-exit() {
  info "waiting for stack to complete"
  for i in {1..30} ; do
    STATUS=$( describe-stack ${KORE_MASTER_ROLE_NAME} StackStatus )

    case "${STATUS}" in
      "CREATE_COMPLETE" | "UPDATE_COMPLETE")
        break
        ;;
      "ROLLBACK_COMPLETE")
        error "Unrecoverable stack status ${STATUS}, please review and delete stack and try again"
        ;;
      "CREATE_FAILED" | "ROLLBACK_FAILED" | "ROLLBACK_IN_PROGRESS")
        error "Stack error status ${STATUS} -  back for ${KORE_MASTER_ROLE_NAME}"
        ;;
      *)
        sleep 1
        ;;
    esac
  done

  if [[ "${STATUS}" =~ ^(CREATE_COMPLETE|UPDATE_COMPLETE)$ ]]; then
    info "Role successfuly created: ${KORE_MASTER_ROLE_NAME}"
  else
    error "Stack didn't complete - ${STATUS}. Reveiw cloudformation stacxk events for ${KORE_MASTER_ROLE_NAME}"
  fi
}

check-dependency() {
  bin=${1?"missing name"}
  which ${bin} >/dev/null 2>&1 || \
    error "missing cli tool:${bin}, please install and retry"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --master-role-name) KORE_MASTER_ROLE_NAME=${2}; shift 2; ;;
    --kore-user-arn)    KORE_USER_ARN=${2};         shift 2; ;;
    --s3-bucket)        S3_BUCKET=${2};             shift 2; ;;
    --s3-region)        S3_REGION=${2};             shift 2; ;;
    --dry-run)          DRY_RUN=true;               shift 1; ;;
    -h|--help)          usage;                               ;;
    -q|-quiet)          QUIET=true;                 shift 1; ;;
    *)                                              shift 1; ;;
  esac
done

check-dependency jq
check-dependency aws

[[ -z ${KORE_USER_ARN:-} ]]                  && usage "You must specify the ARN of the Kore user identity"
[[ -z ${KORE_MASTER_ROLE_NAME:-} ]]          && usage "Please specify master role name (e.g. kore-account-management-role-for-custom-ou)"
[[ -z ${S3_REGION:-} ]]                      && usage "Unknown s3 region, please configure aws default or specify"
[[ ! -f ${MASTER_ROLE_CF_TEMPLATE_PATH} ]]   && usage "Missing file ${MASTER_ROLE_CF_TEMPLATE_NAME}! Did you download it?"

info "stack bucket: ${S3_BUCKET}"
info "stack bucket region: ${S3_REGION}"
info "stack name: ${KORE_MASTER_ROLE_NAME} (to create role of same name)"
info "will grant sts permission to user:${KORE_USER_ARN}"

# These will conflict and result in endpoint error when uploading stack
unset AWS_DEFAULT_REGION AWS_REGION

create-bucket-if-required
deploy-stack
wait-on-stack-complete-or-exit
print-stack-outputs ${KORE_MASTER_ROLE_NAME}
