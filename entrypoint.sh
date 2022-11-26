#!/bin/bash

DEBUG=$8

if [[ "${DEBUG}" == "true" ]]; then
  echo "Current directory: `pwd`"
  echo "Contents: `ls -la`"
  echo ""
  echo "[Settings]"
  echo "schema-id: ${1}"
  echo "schema-type: ${2}"
  echo "schema-name: ${3}"
  echo "descriptor-set: ${4}"
  echo "api-token: ${5}"
fi

/publish -schema-id ${1} \
  -schema-type ${2} \
  -schema-name "${3}" \
  -descriptor-set ${4} \
  -api-token ${5}
