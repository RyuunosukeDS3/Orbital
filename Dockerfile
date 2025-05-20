FROM debian:bookworm-slim

ARG TARGETARCH
ARG APP_NAME=orbital

COPY build/${APP_NAME}-linux-${TARGETARCH} /usr/local/bin/${APP_NAME}
RUN chmod 755 /usr/bin/${APP_NAME}

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/orbital"]
