#!/bin/bash

DEBUG=$8

if [[ "${DEBUG}" == "true" ]]; then
  echo "Current directory: $(pwd)"
  echo "Contents: $(ls -la)"
  echo ""
  echo "[Settings]"
  echo "schema-id: ${STREAMDAL_SCHEMA_ID}"
  echo "schema-type: ${STREAMDAL_SCHEMA_TYPE}"
  echo "schema-name: ${STREAMDAL_SCHEMA_NAME}"
  echo "api-token: ${STREAMDAL_API_TOKEN}"
  echo "api-address: ${STREAMDAL_SCHEMA_API_ADDRESS}"
  echo "input: ${STREAMDAL_INPUT}"
  echo "input type: ${STREAMDAL_INPUT_TYPE}"
  echo "output: ${STREAMDAL_OUTPUT}"
fi

/publish
