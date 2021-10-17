#!/bin/bash

# Required Variables
# -------------------------
# GCP_SERVICE_ID
# GCP_PROJECT_NUMBER
# GCP_PROJECT_ID
# GCP_REGION
# GCP_SERVICE_ACCOUNT
# GCP_DB_NAME
# OCTOVY_URL
# OCTOVY_GITHUB_APP_ID
# OCTOVY_GITHUB_CLIENT_ID
# OCTOVY_DB_TYPE

IMAGE=$1

gcloud beta run deploy "${GCP_SERVICE_ID}" \
    --project="${GCP_PROJECT_ID}" \
    --image "${IMAGE}" \
    --region="${GCP_REGION}" \
    --platform="managed" \
    --cpu=1  \
    --memory=512Mi \
    --port 9080 \
    --args serve \
    --allow-unauthenticated \
    --service-account=${GCP_SERVICE_ACCOUNT} \
    --set-cloudsql-instances="octovy-service:asia-northeast1:${GCP_DB_NAME}" \
    --set-env-vars "OCTOVY_LOG_LEVEL=debug,
    OCTOVY_LOG_FORMAT=json,
    OCTOVY_FRONTEND_URL=${OCTOVY_URL},
    OCTOVY_GITHUB_APP_ID=${OCTOVY_GITHUB_APP_ID},
    OCTOVY_GITHUB_CLIENT_ID=${OCTOVY_GITHUB_CLIENT_ID},
    OCTOVY_ADDR=0.0.0.0,
    OCTOVY_DB_TYPE=${OCTOVY_DB_TYPE},
    GIN_MODE=release" \
    --set-secrets "OCTOVY_DB_CONFIG=projects/${GCP_PROJECT_NUMBER}/secrets/OCTOVY_DB_CONFIG:latest,
    OCTOVY_GITHUB_WEBHOOK_SECRET=projects/${GCP_PROJECT_NUMBER}/secrets/OCTOVY_GITHUB_WEBHOOK_SECRET:latest,
    OCTOVY_GITHUB_APP_PRIVATE_KEY=projects/${GCP_PROJECT_NUMBER}/secrets/OCTOVY_GITHUB_APP_PRIVATE_KEY:latest,
    OCTOVY_GITHUB_SECRET=projects/${GCP_PROJECT_NUMBER}/secrets/OCTOVY_GITHUB_SECRET:latest"
