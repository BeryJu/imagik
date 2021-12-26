# Build WebUI First
FROM --platform=${BUILDPLATFORM} docker.io/node:16 as npm-builder

COPY ./web/ /build/web

RUN cd /build/web && npm i && npm run build

# Build application second
FROM docker.io/golang:1.17.5 AS builder

ENV CGO_ENABLED=0
ARG GIT_BUILD_HASH
ENV GIT_BUILD_HASH=$GIT_BUILD_HASH

COPY . /go/src/beryju.org/imagik
COPY --from=npm-builder /build/root /go/src/beryju.org/imagik/root

RUN cd /go/src/beryju.org/imagik && \
    go build -X main.buildCommit=$GIT_BUILD_HASH -v -o /go/bin/imagik

# Final container
FROM gcr.io/distroless/static-debian11:debug

COPY --from=builder /go/bin/imagik /imagik
COPY ./config.docker.yml /config.yml

EXPOSE 8000

WORKDIR /share

ENV IMAGIK_ROOT=/share
ENV IMAGIK_AUTH_DRIVER=null

ENTRYPOINT [ "/imagik", "-c=/config.yml" ]
