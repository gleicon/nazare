### nazare
Nazaré is a server for sketch counters persisted on badger db with a redis interface.

### Why
Drop in replacement for Redis specifically on counters to achieve better performance using multicore and a specialised persistence mechanism with badger db.

### Build and run

$ make

$ ./nazare

### Options

-s ip:port - ip and port to bind for redis protocol, default 0.0.0.0:6379

-d dbpath - hllcounters.db

-a api ip:port for http api and metrics - default 127.0.0.1:8080

### Implemented commands

	PFADD
	PFCOUNT

### TODO

	implement PFMERGE
	implement GET/SET w/ cuckoo filter
	metrics and stats

### Nazaré

![nazarect](nazare.jpg)

gleicon
