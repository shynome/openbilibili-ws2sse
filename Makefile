VERSION := $(shell git describe --tags --always --dirty)

build:
	CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${VERSION}' -s -w" -o openbilibili-ws2sse .
docker: build
	docker build . -t shynome/openbilibili-ws2sse:${VERSION}
push: docker
	docker push shynome/openbilibili-ws2sse:${VERSION}
