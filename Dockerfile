# Build WebUI First
FROM node as npm-builder

COPY ./web/ /build/web

RUN cd /build/web && npm i && npm run build

# Build application second
FROM golang:latest AS builder

WORKDIR /usr/local/go/src/gopyazo

COPY . /usr/local/go/src/gopyazo
COPY --from=npm-builder /build/root /usr/local/go/src/gopyazo/root

RUN apt-get update && \
    apt-get install -y --no-install-recommends libmagic-dev && \
    rm -rf /var/lib/apt/lists/*

RUN make build

# Final container
FROM debian

RUN apt-get update && \
    apt-get install -y --no-install-recommends libmagic-dev && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/gopyazo /gopyazo
COPY ./config.docker.yml /config.yml

EXPOSE 8000

WORKDIR /share

ENV GOPYAZO_ROOT=/share
ENV GOPYAZO_AUTH_DRIVER=null

ENTRYPOINT [ "/gopyazo", "-c=/config.yml" ]
