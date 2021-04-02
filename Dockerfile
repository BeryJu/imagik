# Build WebUI First
FROM node as npm-builder

COPY ./web/ /build/web

RUN cd /build/web && npm i && npm run build

# Build application second
FROM golang:1.16.3 AS builder

COPY . /go/src/github.com/BeryJu/imagik
COPY --from=npm-builder /build/root /go/src/github.com/BeryJu/imagik/root

RUN cd /go/src/github.com/BeryJu/imagik && make docker-build

# Final container
FROM debian

COPY --from=builder /go/bin/imagik /imagik
COPY ./config.docker.yml /config.yml

EXPOSE 8000

WORKDIR /share

ENV IMAGIK_ROOT=/share
ENV IMAGIK_AUTH_DRIVER=null

ENTRYPOINT [ "/imagik", "-c=/config.yml" ]
