package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/gleicon/nazare/server"
	"github.com/gleicon/nazare/sets"

	"github.com/gleicon/nazare/counters"
	"github.com/gleicon/nazare/db"
)

func help() {
	fmt.Println("hyperloglog counter server")
	fmt.Println("Talks redis protocol and limited commandset")
	fmt.Println("-s ip:port - ip and port to bind for redis protocol, default 0.0.0.0:6379")
	fmt.Println("-d dbpath - hllcounters.db")
	fmt.Println("-m api ip:port for prometheus metrics - default 127.0.0.1:2112")
	os.Exit(-1)
}

var localCounters *counters.HLLCounters
var localSets *sets.CkSet
var localDatastorage db.Datastorage

func main() {

	var serverAddr, httpAPIAddr, dbPath string

	serverAddr = *flag.String("s", "127.0.0.1:6379", "Redis server ip:port")
	httpAPIAddr = *flag.String("m", "127.0.0.1:2112", "Prometheus metrics ip:port")
	dbPath = *flag.String("d", "hllcounters.db", "Database path")

	flag.Usage = help
	flag.Parse()

	nzs, err := server.NewNZServer(serverAddr, httpAPIAddr, dbPath)
	if err != nil {
		log.Panic("Error creating server:", err)
	}
	err = nzs.Start()
	if err != nil {
		log.Panic("Error starting server:", err)
	}
}
