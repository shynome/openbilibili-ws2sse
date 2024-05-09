build:
	CGO_ENABLED=0 go build -ldflags="-X 'main.Version=$$(git describe --tags --always --dirty)' -s -w" -o openbilibili-ws2sse .
docker: build
	docker build . -t shynome/openbilibili-ws2sse:$$(git describe --tags --always --dirty)
push: docker
	docker push shynome/openbilibili-ws2sse:$$(git describe --tags --always --dirty)
