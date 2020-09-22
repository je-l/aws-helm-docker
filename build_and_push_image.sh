#!/usr/bin/env bash

set -euxo pipefail


CREATED_DATE=$(date --rfc-3339=date)
export $(go run fetch_versions.go)
IMAGE_NAME=jljljl/aws-helm:$CREATED_DATE

docker build \
  --build-arg CREATED_DATE="$CREATED_DATE" \
  --build-arg KUBECTL_VERSION="$KUBECTL_VERSION" \
  --build-arg HELM_VERSION="$HELM_VERSION" \
  --build-arg AWS_CLI_VERSION="$AWS_CLI_VERSION" \
  -t "$IMAGE_NAME" \
  .

docker push $IMAGE_NAME
