# Build WebUI First
FROM --platform=${BUILDPLATFORM} docker.io/node:16 as npm-builder

COPY ./web/ /build/web

RUN cd /build/web && npm i && npm run build

# Build application second
FROM docker.io/golang:1.17.5 AS builder

COPY . /go/src/github.com/BeryJu/imagik
COPY --from=npm-builder /build/root /go/src/github.com/BeryJu/imagik/root

RUN cd /go/src/github.com/BeryJu/imagik && make docker-build

# Final container
FROM docker.io/debian

COPY --from=builder /go/bin/imagik /imagik
COPY ./config.docker.yml /config.yml

EXPOSE 8000

WORKDIR /share

ENV IMAGIK_ROOT=/share
ENV IMAGIK_AUTH_DRIVER=null

ENTRYPOINT [ "/imagik", "-c=/config.yml" ]
