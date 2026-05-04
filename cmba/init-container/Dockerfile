# cmba-verify init container image.
# Smallest possible: busybox + yq for YAML field extraction.
FROM busybox:1.36-musl

# Install yq (Go binary for YAML query).
ADD https://github.com/mikefarah/yq/releases/download/v4.40.5/yq_linux_amd64 /usr/local/bin/yq
RUN chmod 0755 /usr/local/bin/yq

COPY verify.sh /usr/local/bin/cmba-verify
RUN chmod 0755 /usr/local/bin/cmba-verify

USER 65532:65532
ENTRYPOINT ["/usr/local/bin/cmba-verify"]
