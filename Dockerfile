FROM alpine:3.21
COPY bin/metal-linux-amd64 /metal
ENTRYPOINT ["/metal"]
