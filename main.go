package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gleicon/renand/counters"
	"github.com/gleicon/renand/db"
)

func help() {
	fmt.Println("increment hll counters: counters -n <countername> -i foobar")
	fmt.Println("estimate counter size: counters -n <countername>")
	os.Exit(-1)
}

func main() {

	counterName := flag.String("n", "", "HLL counter name")
	incrementBy := flag.String("i", "", "Increment counter by string")
	badger := flag.Bool("b", false, "Enable badger datastorage")

	flag.Usage = help
	flag.Parse()

	var counters *counters.HLLCounters

	if *badger == true {
		badgerds, _ := db.NewBadgerDatastorage("testcounters")
		counters, _ = counters.NewHLLCounters(badgerds)
	} else {
		memds, _ := db.NewHLLDatastorage()
		counters, _ = counters.NewHLLCounters(memds)

	}

	if *counterName == "" {
		fmt.Println("No counter")
		os.Exit(-1)
	}
	var cc uint64
	var err error

	if *incrementBy == "" {
		if cc, err = counters.RetrieveCounterEstimate(*counterName); err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)

		}
		fmt.Printf("counter: %s - est. size: %d\n", *counterName, cc)
	} else {
		if err := counters.IncrementCounter(*counterName, []byte(*incrementBy)); err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		if cc, err = counters.RetrieveCounterEstimate(*counterName); err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)

		}
		fmt.Printf("counter: %s - est. size: %d\n", *counterName, cc)
		xx, _ := counters.ExportCounter(*counterName)
		fmt.Printf("%x\n", xx)

	}

}
