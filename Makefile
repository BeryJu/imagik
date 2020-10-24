run:
ifeq (,$(wildcard ./config.local.yml))
	cp config.yml config.local.yml
endif
	go run -v . server -c config.local.yml
