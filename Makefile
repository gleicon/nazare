include Makefile.defs

all: deps server cli

deps:
	make -C db deps                                                        
	make -C counters deps 
	make -C sets deps 
	make -C nazare-cli deps 
	go get -v

test:
	make -C db test
	make -C counters test
	make -C sets test
	make -C nazare-cli test
	go test -v

server:
	go build -v -o $(NAME) 

cli:
	make -C all nazare-cli

clean:
	rm -f $(NAME)
	make -C nazare-cli clean

.PHONY: server
