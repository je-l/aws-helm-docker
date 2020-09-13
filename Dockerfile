ARG AWS_CLI_VERSION

FROM amazon/aws-cli:${AWS_CLI_VERSION}

ARG AWS_CLI_VERSION
ARG CREATED_DATE
ARG HELM_VERSION
ARG KUBECTL_VERSION

RUN curl -LSs "https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" > /usr/bin/kubectl && chmod +x /usr/bin/kubectl
RUN curl -LSs "https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz" > helm.tar.gz && \
    yum install -y tar gzip && \
    tar xvf helm.tar.gz && \
    cp linux-amd64/helm /usr/bin/helm && \
    rm -r helm*.tar.gz linux-amd64 /var/cache/yum && \
    yum remove -y tar gzip && \
    yum clean all

# https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.created="${CREATED_DATE}"
LABEL org.opencontainers.image.url="https://github.com/je-l/aws-helm-docker"
LABEL org.opencontainers.image.description="CLI tools for AWS EKS: Helm ${HELM_VERSION}, aws-cli ${AWS_CLI_VERSION} and kubectl ${KUBECTL_VERSION}"

ENTRYPOINT ["/bin/sh"]
