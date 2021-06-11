.PHONY: build install

build:
	go build -o build/ots .

install: build
	rm -f "/usr/local/bin/ots"
	ln -s "$(shell pwd)/build/ots" /usr/local/bin/ots

