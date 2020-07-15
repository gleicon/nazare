### nazare
Nazaré is a server for sketch counters and sets persisted on badger db with a redis interface. It is a database that is not so sure about its data.

### Why
Many opensource services resort to sketch structures when the cardinality (size) of a counter reaches beyond a certain value. [Elasticsearch](https://www.elastic.co/blog/count-elasticsearch) uses hyperloglog for counters over 40k documents. Cassandra uses bloom filter to prevent disk hits as cache. Redis implements HyperLogLog too. 

Sketch structures trade size or performance by accuracy. Different implementations are available that tune these parameters. 

Nazare is a drop in replacement for Redis as it speaks the same protocol, enabling any application that implements a Redis Driver to use such counters and sets operations.

### Database
The Underlying database is [BadgerDB](https://github.com/dgraph-io/badger), which implements a series of improvements over non Golang local kv values, including concurrent ACID transactions.

Having a local database is counterintuitive if you building a distributed system but Nazare is meant as a lightweight server and command line so you can `cp` or `rsync` your database around servers, download from S3 before start and delegate loading it to another process, effectively rendering nazare as a cache for sketch structures on ephemeral servers.


### Build and run

$ make

$ ./nazare

### Options

-s ip:port - ip and port to bind for redis protocol, default 0.0.0.0:6379

-d dbpath - hllcounters.db

-a api ip:port for http api and metrics - default 127.0.0.1:8080

### Implemented commands

	GET
	SET
	DEL
	PFADD
	PFCOUNT
	SADD
	SREM
	SCARD
	SISMEMBER

### TODO
	hot database reload

### Nazaré

![nazarect](nazare.jpg)

gleicon
