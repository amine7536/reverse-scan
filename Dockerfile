FROM alpine:latest

COPY reverse-scan /usr/local/bin/reverse-scan

ENTRYPOINT ["/usr/local/bin/reverse-scan"]
