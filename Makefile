run:
ifeq (,$(wildcard ./config.local.yml))
	cp config.yml config.local.yml
endif
	go run -v . -c config.local.yml

docker:
	docker build -t ghcr.io/beryju/imagik:latest .
	docker run -p 8000:8000 -it --rm ghcr.io/beryju/imagik:latest
