# Build WebUI First
FROM node as npm-builder

COPY ./web/ /build/web

RUN cd /build/web && npm i && npm run build

# Build application second
FROM golang:1.15 AS builder

COPY . /go/src/github.com/BeryJu/gopyazo
COPY --from=npm-builder /build/root /go/src/github.com/BeryJu/gopyazo/root

RUN cd /go/src/github.com/BeryJu/gopyazo && make docker-build

# Final container
FROM debian

COPY --from=builder /go/bin/gopyazo /gopyazo
COPY ./config.docker.yml /config.yml

EXPOSE 8000

WORKDIR /share

ENV GOPYAZO_ROOT=/share
ENV GOPYAZO_AUTH_DRIVER=null

ENTRYPOINT [ "/gopyazo", "-c=/config.yml" ]
