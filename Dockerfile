FROM golang:1.12
ADD . /app
WORKDIR /app
RUN make install
ENTRYPOINT [ "/go/bin/reverse-scan" ]