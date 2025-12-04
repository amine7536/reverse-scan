FROM golang:1.23
ADD . /app
WORKDIR /app
RUN make install
ENTRYPOINT [ "/go/bin/reverse-scan" ]