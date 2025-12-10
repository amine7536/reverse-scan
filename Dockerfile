FROM golang:1.25.5
ADD . /app
WORKDIR /app
RUN make install
ENTRYPOINT [ "/go/bin/reverse-scan" ]