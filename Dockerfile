FROM golang:1.25.4
ADD . /app
WORKDIR /app
RUN make install
ENTRYPOINT [ "/go/bin/reverse-scan" ]