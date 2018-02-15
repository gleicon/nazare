package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gleicon/nazare/counters"
	"github.com/gleicon/nazare/db"
	"github.com/tidwall/redcon"
)

func help() {
	fmt.Println("hyperloglog counter server")
	fmt.Println("Talks redis protocol and limited commandset")
	fmt.Println("-s ip:port - ip and port to bind for redis protocol, default 0.0.0.0:6379")
	fmt.Println("-d dbpath - hllcounters.db")
	fmt.Println("-a api ip:port for http api and metrics - default 127.0.0.1:8080")
	os.Exit(-1)
}

var localCounters *counters.HLLCounters

func main() {

	var serverAddr, httpAPIAddr, dbPath string

	serverAddr = *flag.String("s", "0.0.0.0:6379", "Redis server ip:port")
	httpAPIAddr = *flag.String("a", "127.0.0.1:8080", "Api and metrics ip:port")
	dbPath = *flag.String("d", "hllcounters.db", "Database path")

	flag.Usage = help
	flag.Parse()

	badgerds, _ := db.NewBadgerDatastorage(dbPath)
	localCounters, _ = counters.NewHLLCounters(badgerds)
	log.Println(httpAPIAddr)

	go func() {
		err := redcon.ListenAndServe(serverAddr, redisCommandParser, newConnection, closeConnection)
		log.Fatal(err)
	}()
	select {}
}
