#!/bin/bash

pwd
ls

echo "schema-id: ${1}"
echo "schema-type: ${2}"
echo "schema-name: ${3}"
echo "api-token: ${4}"
echo "root-dir: ${5}"
echo "root-message: ${6}"
echo "api-address: ${7}"

/publish -schema-id $1 \
  -schema-type $2 \
  -schema-name "${3}" \
  -api-token $4 \
  -root-dir "$5" \
  -root-message "$6" \
  -api-address "${7}"
