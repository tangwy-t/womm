FROM alpine:3.19
RUN apk add --no-cache ca-certificates && mkdir -p /data
COPY womm /usr/local/bin/womm
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["womm", "serve", "-c", "/data/womm.toml"]
