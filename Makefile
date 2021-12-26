run:
	packr2
ifeq (,$(wildcard ./config.local.yml))
	cp config.yml config.local.yml
endif
	go run -v . -c config.local.yml

docker:
	packr2 clean
	docker build -t ghcr.io/beryju/imagik:latest .
	docker run -p 8000:8000 -it --rm ghcr.io/beryju/imagik:latest

docker-build:
	go get -u github.com/gobuffalo/packr/v2/packr2
	packr2
	go build -v -o /go/bin/imagik
