#!/bin/bash

if [[ "${DEBUG}" == "true" ]]; then
  echo "Current directory: $(pwd)"
  echo "Contents: $(ls -la)"
  echo ""
  echo "[Settings]"
  echo "schema-id: ${BATCH_SCHEMA_ID}"
  echo "schema-type: ${BATCH_SCHEMA_TYPE}"
  echo "schema-name: ${BATCH_SCHEMA_NAME}"
  echo "api-token: ${BATCH_API_TOKEN}"
  echo "root-dir: ${BATCH_ROOT_DIR}"
  echo "api-address: ${BATCH_SCHEMA_API_ADDRESS}"
  echo "artifact type: ${BATCH_ARTIFACT_TYPE}"
  echo "descriptor set path: ${BATCH_DESCRIPTOR_SET_PATH}"
fi

/publish