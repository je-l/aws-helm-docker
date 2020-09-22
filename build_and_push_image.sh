#!/usr/bin/env bash

set -euxo pipefail

export $(go run fetch_versions.go)

CREATED_DATE=$(date --rfc-3339=date)
NEW_IMAGE=jljljl/aws-helm:$CREATED_DATE
LATEST_IMAGE=jljljl/aws-helm:latest

docker build \
  --build-arg CREATED_DATE="$CREATED_DATE" \
  --build-arg KUBECTL_VERSION="$KUBECTL_VERSION" \
  --build-arg HELM_VERSION="$HELM_VERSION" \
  --build-arg AWS_CLI_VERSION="$AWS_CLI_VERSION" \
  -t "$NEW_IMAGE" \
  .

docker tag "$NEW_IMAGE" "$LATEST_IMAGE"
docker push "$NEW_IMAGE"
docker push "$LATEST_IMAGE"
