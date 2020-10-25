run:
	packr2
ifeq (,$(wildcard ./config.local.yml))
	cp config.yml config.local.yml
endif
	go run -v . server -c config.local.yml

docker:
	packr2 clean
	docker build -t beryju/gopyazo:latest .
	docker run -p 8000:8000 -it --rm beryju/gopyazo:latest server

docker-build:
	go get -u github.com/gobuffalo/packr/v2/packr2
	packr2
	go build -v -o /go/bin/gopyazo
