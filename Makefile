all: reverse-scan

GO ?= go
GOTEST = go test -v -bench\=.
WORKDIR ?= $(shell pwd)

reverse-scan:
	mkdir -p build
	$(GO) env
	$(GO) build -ldflags="-s -w" $(EXTRA_BUILD_FLAGS) -o build/reverse-scan

clean:
	rm -f build/reverse-scan

image:
	docker build -t reverse-scan .

build: image
	docker run -v $(WORKDIR):/go/src/github.com/amine7536/reverse-scan -it reverse-scan