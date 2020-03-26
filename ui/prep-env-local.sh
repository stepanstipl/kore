#!/bin/bash

if [ ! -f ../demo.env ] ; then
    echo "You could copy and edit the file:"
    echo "    cp ../hack/compose/demo.env.tmpl ../demo.env"
    echo "... or run bin/korectl local configure"
    exit 1
fi
source ../demo.env
export \
    KORE_IDP_CLIENT_ID \
    KORE_IDP_CLIENT_SECRET \
    KORE_IDP_SERVER_URL \
    KORE_IDP_USER_CLAIMS \
    KORE_IDP_CLIENT_SCOPES

npm run dev