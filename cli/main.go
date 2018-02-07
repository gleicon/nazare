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

	var localCounters *counters.HLLCounters

	if *badger == true {
		badgerds, _ := db.NewBadgerDatastorage("testcounters")
		localCounters, _ = counters.NewHLLCounters(badgerds)
	} else {
		memds, _ := db.NewHLLDatastorage()
		localCounters, _ = counters.NewHLLCounters(memds)

	}

	if *counterName == "" {
		fmt.Println("No counter")
		os.Exit(-1)
	}
	var cc uint64
	var err error

	if *incrementBy == "" {
		if cc, err = localCounters.RetrieveCounterEstimate(*counterName); err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)

		}
		fmt.Printf("counter: %s - est. size: %d\n", *counterName, cc)
	} else {
		if err := localCounters.IncrementCounter(*counterName, []byte(*incrementBy)); err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		if cc, err = localCounters.RetrieveCounterEstimate(*counterName); err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)

		}
		fmt.Printf("counter: %s - est. size: %d\n", *counterName, cc)
		xx, _ := localCounters.ExportCounter(*counterName)
		fmt.Printf("%x\n", xx)

	}

}
