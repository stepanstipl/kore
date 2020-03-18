#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

address=${1?'expected address'}
for i in {1..10} ; do
    ret=0
    curl \
        --retry 50 \
        --retry-delay 3 \
        --retry-connrefused \
        -sSL ${address} >/dev/null 2>&1 || ret=$?
    if [ ${ret} -eq 52 ]; then
        echo "empty response with:$1, trying again..."
        sleep 5
        continue
    fi
    if [ ${ret} -ne 0 ]; then break ; fi
done
exit $ret
