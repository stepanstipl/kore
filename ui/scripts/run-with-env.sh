#!/bin/sh

set -euo pipefail

for idpvar in KORE_IDP_CLIENT_ID KORE_IDP_CLIENT_SECRET KORE_IDP_SERVER_URL KORE_IDP_USER_CLAIMS KORE_IDP_CLIENT_SCOPES; do
  export ${idpvar}=$(kubectl --context kind-kore -n kore get secret kore-idp -o json | jq -r ".data.${idpvar}" | base64 --decode)
done

export KORE_API_TOKEN=$(kubectl --context kind-kore -n kore get secret kore-api -o json | jq -r ".data.KORE_ADMIN_TOKEN" | base64 --decode)

export KORE_FEATURE_GATES="services=true,application_services=true,monitoring_services=true"
export KORE_UI_SHOW_PROTOTYPES=true

exec "$@"
