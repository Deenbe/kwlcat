.PHONY: build install

build:
	go build -o build/gopi .

install: build
	rm -f "/usr/local/bin/gopi"
	ln -s "$(shell pwd)/build/gopi" /usr/local/bin/gopi

