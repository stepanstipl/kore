#!/bin/sh

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

${DIR}/run-with-env.sh npm start
