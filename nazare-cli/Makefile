NAME=nazare-cli
VERSION=0.0.1

all: deps cli

deps:
	go get -v

test:
	go test -v

cli:
	go build -v -o $(NAME) 

clean:
	rm -f $(NAME)

.PHONY: server
