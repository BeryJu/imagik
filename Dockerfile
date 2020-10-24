FROM golang:latest AS builder

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -d -v ./...

RUN go build -v -o /go/bin/gopyazo

FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/bin/gopyazo /gopyazo
COPY ./config.docker.yml /config.yml
EXPOSE 8000
WORKDIR /share
ENV GOPYAZO_ROOT=/share
ENV GOPYAZO_AUTH_DRIVER=null
ENTRYPOINT [ "/gopyazo", "-c=/config.yml" ]
