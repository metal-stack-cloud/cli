FROM alpine:3.23
COPY bin/metal-linux-amd64 /metal
ENTRYPOINT ["/metal"]
