FROM alpine:latest

ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/reverse-scan /usr/local/bin/reverse-scan

ENTRYPOINT ["/usr/local/bin/reverse-scan"]
