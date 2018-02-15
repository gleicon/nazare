all: deps cli

deps:
	go get -v

test:
	make -C db test
	make -C counters test
	go test -v

cli:
	go build -v -o nazare-local

clean:
	rm -f nazare-local

.PHONY: server
