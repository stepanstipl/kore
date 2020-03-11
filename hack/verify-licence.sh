#!/bin/bash
#
# Description: is used to check all the files have the same headers
#

BOILERPLATE=${BOILERPLATE:-"hack/boilerplate.go.txt"}
BOILERPLATE_LENGTH=$(cat ${BOILERPLATE}| wc -l | xargs)
EXCLUDE_FILES=(
  ./hack/generate/manifests_vfsdata.go
  ./pkg/clusterappman/manifests_tools.go
  ./pkg/clusterappman/manifests_vfsdata.go
  ./pkg/tools/tools.go
)

if [[ -z "${BOILERPLATE_LENGTH}" ]]; then
  echo "Failed to retrieve length of header in ${BOIERPLATE}"
  exit 1
fi

while read name; do
  # ignore excluded files
  [[ " ${EXCLUDE_FILES[*]} " == *" ${name} "* ]] && continue
  # ignore auto generated ones
  [[ "${name}" =~ ^.*zz_generated.*$ ]] && continue

  if ! head -n ${BOILERPLATE_LENGTH} ${name} | diff - ${BOILERPLATE} >/dev/null; then
    echo "Please check the licence header on ${name}"
    echo "Ensure its the same as ${BOILERPLATE}"
  fi
done < <(find . -type f -name "*.go" | grep -v vendor)
