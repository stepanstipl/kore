#!/bin/bash -eu
#
# Copyright (C) 2020 Appvia Ltd <info@appvia.io>
#
# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation; either version 2
# of the License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#

## Set the defaults
BUILD_AUTH_PROXY=true
BUILD_KORE_API=true
BUILD_CLI=true
BUILD_IMAGES=true
ENABLE_API=true
ENABLE_UI=false
K8S_VERSION=""
VERSION=""

# Make this pretty
export NC='\e[0m'
export GREEN='\e[0;32m'
export YELLOW='\e[0;33m'
export RED='\e[0;31m'
export PATH=${PATH}:${PWD}/bin

log()      { (2>/dev/null echo -e "$@${NC}"); }
announce() { log "${GREEN}[$(date +"%T")] [INFO] $@"; }
failed()   { log "${YELLOW}[$(date +"%T")] [FAIL] $@"; }

usage() {
  cat <<EOF
  Usage: $(basename $0)
  --build-kore-api <bool>  : indicates we should build the kore-apiserver ourselves (defaults: ${BUILD_KORE_API})
  --build-cli    <bool>    : indicates should should build the kore cli (defaults: ${BUILD_CLI})
  --build-images <bool>    : indicates if we should build the images locally (defaults: ${BUILD_IMAGES})
  --build-proxy  <bool>    : indicates if we should build the auth proxy image (defaults: ${BUILD_AUTH_PROXY})
  --enable-api   <bool>    : enable the kore-apiserver deployment (defaults: ${ENABLE_API})
  --enable-ui    <bool>    : enable building and deploying the ui (defaults: ${ENABLE_UI})
  --k8s-version  <string>  : the version of kubernetes we should run against (defaults: ${K8S_VERSION})
  --version      <string>  : is the version name to build the components (default: "")
  -h|--help                : display this usage menu
EOF
  if [[ -n $@ ]]; then
    echo "[error] $@"
    exit 1
  fi
  exit 0
}

build-cli() {
  if [[ ${BUILD_CLI} == true ]]; then
    announce "Building the Kore CLI"
    make kore
  fi
}

build-cluster() {
  announce "Provisioning the local kubernetes cluster"
  local args=""

  if [[ ${ENABLE_UI} == true ]]; then
    announce "Deploying the Kore ui"
    args="${args} --set=ui.hostPort=3000"
    args="${args} --set=ui.replicas=01"
    args="${args} --set=ui.serviceType=NodePort"
    args="${args} --set=ui.disable_animations=true"

    if [[ ${BUILD_IMAGES} == true ]]; then
      announce "Building the Kore UI Image"
      args="${args} --kind-load-image=quay.io/appvia/kore-ui:${VERSION}"
      cd ui
      VERSION=${VERSION} make docker
      cd ..
    fi
  else
    args="${args} --disable-ui"
    args="${args} --set=ui.replicas=00"
  fi

  if [[ ${ENABLE_API} == true ]]; then
    announce "Deploying the Kore API"
    args="${args} --set=api.hostPort=10080"
    args="${args} --set=api.replicas=01"
    args="${args} --set=api.serviceType=NodePort"

    if [[ ${BUILD_IMAGES} == true ]]; then
      args="${args} --kind-load-image=quay.io/appvia/kore-apiserver:${VERSION}"
      # @step: do we need to compile the binary our self?
      if [[ ${BUILD_KORE_API} == true ]]; then
        VERSION=${VERSION} make kore-apiserver
      fi
      VERSION=${VERSION} make kore-apiserver-image-local
      # we make and we push the images as it have to be remote
      if [[ ${BUILD_AUTH_PROXY} == true ]]; then
        announce "Building & pushing the auth proxy image"

        args="${args} --set=api.images={}"
        args="${args} --set=api.images.auth_proxy=quay.io/appvia/auth-proxy:${VERSION}"
        VERSION=${VERSION} make auth-proxy-image-release
      fi
    fi
  else
    args="${args} --set=api.replicas=00"
  fi

  if [[ -n "${K8S_VERSION}" && "${K8S_VERSION}" != "latest" ]]; then
    announce "Using kubenetes version: ${K8S_VERSION}"
    args="${args} --kind-image=kindest/node:${K8S_VERSION}"
  fi

  kore alpha local up \
    --deployment-timeout=8m \
    --force=true \
    --release=charts/kore \
    --set=api.admin_pass=${KORE_ADMIN_PASS} \
    --set=api.admin_token=${KORE_ADMIN_TOKEN} \
    --set=api.auth_plugin_config={} \
    --set=api.auth_plugin_config.local_jwt_publickey=${KORE_LOCAL_JWT_PUBLIC_KEY} \
    --set=api.auth_plugins.0=basicauth \
    --set=api.auth_plugins.1=admintoken \
    --set=api.auth_plugins.2=localjwt \
    --set=api.auth_plugins.3=openid \
    --set=api.verbose=true \
    --set=idp.client_id=${KORE_IDP_CLIENT_ID} \
    --set=idp.client_secret=${KORE_IDP_CLIENT_SECRET} \
    --set=idp.server_url=${KORE_IDP_SERVER_URL} \
    --version=${VERSION} \
    ${args}
}

while [[ $# -gt 0 ]]; do
  case "$1" in
  --build-kore-api) BUILD_KORE_API=${2};   shift 2; ;;
  --build-images)   BUILD_IMAGES=${2};     shift 2; ;;
  --build-proxy)    BUILD_AUTH_PROXY=${2}; shift 2; ;;
  --build-cli)      BUILD_CLI=${2};        shift 2; ;;
  --enable-api)     ENABLE_API=$2;         shift 2; ;;
  --enable-ui)      ENABLE_UI=${2};        shift 2; ;;
  --k8s-version)    K8S_VERSION=${2};      shift 2; ;;
  --version)        VERSION=${2};          shift 2; ;;
  -h|--help)        usage;                          ;;
  *)                                       shift 1; ;;
  esac
done

if [[ -z "${VERSION}" ]]; then
  failed "no version has been set"
  exit 1
fi

build-cli
build-cluster
