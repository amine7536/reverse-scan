FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY reverse-scan /usr/local/bin/reverse-scan

ENTRYPOINT ["/usr/local/bin/reverse-scan"]
