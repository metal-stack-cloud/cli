FROM alpine:3.24
COPY bin/metal-linux-amd64 /metal
ENTRYPOINT ["/metal"]
