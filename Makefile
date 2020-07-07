include Makefile.defs

all: deps server

deps:
	make -C db deps                                                        
	make -C counters deps 
	make -C sets deps 
	go get -v

test:
	make -C db test
	make -C counters test
	make -C sets test
	go test -v

server:
	go build -v -o $(NAME) 

clean:
	rm -f $(NAME)

.PHONY: server
