FROM golang:latest AS builder
WORKDIR $GOPATH/src/github.com/BeryJu/gopyazo
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -v -o /go/bin/gopyazo

FROM scratch
COPY --from=builder /go/bin/gopyazo /gopyazo
EXPOSE 8080
WORKDIR /web-root
CMD [ "/gopyazo" ]
ENTRYPOINT [ "/gopyazo" ]
