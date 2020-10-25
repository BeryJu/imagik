run:
	packr2
ifeq (,$(wildcard ./config.local.yml))
	cp config.yml config.local.yml
endif
	go run -v . server -c config.local.yml

docker:
	docker build -t beryju/gopyazo .
	docker run -it --rm beryju/gopyazo server

build:
	go get -u github.com/gobuffalo/packr/v2/packr2
	packr2
	go build -v -o /go/bin/gopyazo
