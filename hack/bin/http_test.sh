#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

address=${1?'expected address'}
ret=0
for i in 1..3 ; do
    curl \
        --retry 50 \
        --retry-delay 3 \
        --retry-connrefused \
        -sSL ${address} >/dev/null 2>&1 || ret=$?
    if [ ${ret} -eq 52 ]; then
        echo "empty respone, try $1"
        sleep 3
        continue
    fi
    if [ ${ret} -ne 0 ]; then break ; fi
done
exit $ret
