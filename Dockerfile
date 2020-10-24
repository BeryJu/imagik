FROM golang:latest AS builder

WORKDIR /go/src/app
COPY . /go/src/app

RUN apt-get update && \
    apt-get install -y --no-install-recommends libmagic-dev && \
    rm -rf /var/lib/apt/lists/*

RUN go get -d -v ./...

RUN go build -v -o /go/bin/gopyazo

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
