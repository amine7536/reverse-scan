FROM centos:centos7

RUN yum -y install epel-release \
    && yum -y install htop vim wget \
    && yum -y group install "Development Tools"

ENV GOLANG_VERSION 1.8.3

RUN wget -q -O go.tar.gz https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go.tar.gz \
    &&     export PATH="/usr/local/go/bin:$PATH" \
    && go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin"
RUN curl -s  https://glide.sh/get | sh

COPY . $GOPATH/src/github.com/amine7536/reverse-scan
WORKDIR $GOPATH/src/github.com/amine7536/reverse-scan
RUN glide install
ENTRYPOINT  ["make"]
