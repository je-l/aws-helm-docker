FROM amazon/aws-cli:2.0.34

RUN yum install -y wget tar gzip

RUN wget --quiet "https://storage.googleapis.com/kubernetes-release/release/v1.17.9/bin/linux/amd64/kubectl" -O /usr/bin/kubectl \
    && chmod +x /usr/bin/kubectl

RUN wget --quiet "https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz" \
  && tar xvf helm*.tar.gz \
  && cp linux-amd64/helm /usr/bin/helm \
  && rm -r helm*.tar.gz linux-amd64

ENTRYPOINT ["/bin/sh"]
