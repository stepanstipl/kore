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

: ${BUILD_ID?"The build id must be set"}
: ${ENVIRONMENT?"The environment must be set"}
: ${KORE_API_PUBLIC_URL?"The Kore API URL must be set"}
: ${KORE_UI_PUBLIC_URL?"The Kore UI URL must be set"}
: ${VERSION?"The image version must be set"}

CHART=${CHART:-"./charts/kore"}
HELM_STABLE="https://kubernetes-charts.storage.googleapis.com"

log()   { (2>/dev/null echo -e "$@"); }
info()  { log "[info]  $@"; }
error() { log "[error] $@"; exit 1; }

info "adding the helm stable repo"
helm repo add stable ${HELM_STABLE} >/dev/null || error "failed to install the stable repo"

info "updating the helm repos"
helm repo update >/dev/null || error "failed to update the helm repos"

info "deploying the ingress controller to ${ENVIRONMENT} environment"

cat <<EOF > values.yaml
api:
  ingress:
    annotations:
      kubernetes.io/ingress.class: nginx
      nginx.ingress.kubernetes.io/backend-protocol: HTTP
      nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
  images:
    auth_proxy: quay.io/appvia/auth-proxy:${VERSION}
    clusterappman: quay.io/appvia/kore-apiserver:${VERSION}
ui:
  ingress:
    annotations:
      kubernetes.io/ingress.class: nginx
      nginx.ingress.kubernetes.io/backend-protocol: HTTP
      nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
EOF

trap "rm -f values.yaml >/dev/null 2>&1" EXIT

if ! helm upgrade ingress stable/nginx-ingress \
  --install \
  --namespace=ingress \
  --set=controller.autoscaling.enabled=true \
  --set=controller.autoscaling.maxReplicas=10 \
  --set=controller.autoscaling.minReplicas=2 \
  --set=controller.image.tag=0.30.0 \
  --set=controller.service.externalTrafficPolicy=Local \
  --set=controller.service.type=LoadBalancer \
  --set=podSecurityPolicy.enabled=true >/dev/null; then

  error "failed to upgrade or install ingress controller"
fi

info "deploying the version ${VERSION} to ${ENVIRONMENT}, build-id: ${BUILD_ID}"

if ! helm upgrade kore ${CHART} \
  --install \
  --namespace=kore \
  --set=api.build=${VERSION} \
  --set=api.endpoint.url=${KORE_API_PUBLIC_URL} \
  --set=api.ingress.enabled=true \
  --set=api.ingress.hostname=${KORE_API_PUBLIC_URL##https://} \
  --set=api.ingress.tls_secret=kore-app-ingress \
  --set=api.replicas=2 \
  --set=api.version=latest \
  --set=idp.client_id=${KORE_IDP_CLIENT_ID} \
  --set=idp.client_secret=${KORE_IDP_CLIENT_SECRET} \
  --set=idp.server_url=${KORE_IDP_SERVER_URL} \
  --set=ui.build=${VERSION} \
  --set=ui.show_prototypes=${KORE_UI_SHOW_PROTOTYPES:-false} \
  --set=ui.endpoint.url=${KORE_UI_PUBLIC_URL} \
  --set=ui.ingress.enabled=true \
  --set=ui.ingress.hostname=${KORE_UI_PUBLIC_URL##https://} \
  --set=ui.ingress.tls_secret=kore-app-ingress \
  --set=ui.replicas=2 \
  --set=ui.version=${VERSION} \
  --values=values.yaml \
  --wait >/dev/null; then

  error "failed to upgrade or install latest on ${ENVIRONMENT} environment"
fi

info "kore has been upgraded in the ${ENVIRONMENT} environment"

info "checking the api is running "
if ! curl -sL \
  --retry-connrefused \
  --retry 20 \
  --retry-delay 5 \
  ${KORE_API_PUBLIC_URL}/healthz >/dev/null; then

  error "kore api does not appear to be running"
fi
