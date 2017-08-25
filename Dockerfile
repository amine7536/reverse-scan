FROM ubuntu:latest

RUN apt-get -y update \
    && apt-get -y install make curl wget git \
    && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.9

RUN wget -qO go${GOLANG_VERSION}.linux-amd64.tar.gz https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz \
    && export PATH="/usr/local/go/bin:$PATH" \
    && go version \
    && rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin"
RUN wget -qO - https://glide.sh/get | bash

COPY . $GOPATH/src/github.com/amine7536/reverse-scan
WORKDIR $GOPATH/src/github.com/amine7536/reverse-scan
RUN glide install
ENTRYPOINT  ["make"]
