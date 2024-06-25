FROM golang:1.22-alpine AS build

RUN apk add --no-cache bash git make build-base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make bin/reverse-scan

FROM alpine:3 AS runtime

COPY --from=build /build/bin/reverse-scan /usr/local/bin/reverse-scan

ENTRYPOINT [ "/usr/local/bin/reverse-scan" ]
