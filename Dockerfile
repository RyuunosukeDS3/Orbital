FROM debian:bookworm-slim

ARG TARGETARCH
ARG APP_NAME=orbital

COPY build/${APP_NAME}-linux-${TARGETARCH} /usr/local/bin/${APP_NAME}

ENTRYPOINT ["/usr/local/bin/orbital"]
