#!/bin/bash

GIT=`which git`
SWAGGER="swagger.json"

if [[ -z "${GIT}" ]]; then
  echo "[error] failed to find git in PATH ($PATH)"
  exit 1
fi

echo "--> Checking if the ${SWAGGER} has changed"
if git status | grep -q ${SWAGGER}; then
  echo "The API has been updated and the ${SWAGGER} is out of date, please run $ make swagger"
  exit 1
fi
