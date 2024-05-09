FROM alpine:latest
WORKDIR /app
EXPOSE 7070

COPY openbilibili-ws2sse /openbilibili-ws2sse
ENTRYPOINT [ "/openbilibili-ws2sse"]
CMD []
