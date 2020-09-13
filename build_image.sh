#!/usr/bin/env bash

set -euxo pipefail


CREATED_DATE=$(date --rfc-3339=date)
KUBECTL_VERSION=$(curl -LSs https://api.github.com/repos/kubernetes/kubernetes/releases | jq -r '.[] | .tag_name' | grep -E '^v1\.17\.[0-9]+$' | sort --version-sort | tail -n1)
HELM_VERSION=$(curl -LSs https://api.github.com/repos/helm/helm/releases | jq -r  '.[]|.tag_name' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort --version-sort | tail -n1)
AWS_CLI_VERSION=$(curl -LSs https://hub.docker.com/v2/repositories/amazon/aws-cli/tags/ | jq -r '.results[] | .name' | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$' | sort --version-sort | tail -n1)
IMAGE_NAME=jljljl/aws-helm:$CREATED_DATE

docker build \
  --build-arg CREATED_DATE="$CREATED_DATE" \
  --build-arg KUBECTL_VERSION="$KUBECTL_VERSION" \
  --build-arg HELM_VERSION="$HELM_VERSION" \
  --build-arg AWS_CLI_VERSION="$AWS_CLI_VERSION" \
  -t "$IMAGE_NAME" \
  .
