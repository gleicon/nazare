package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gleicon/nazare/counters"
	"github.com/gleicon/nazare/db"
)

func help() {
	fmt.Println("nazare-cli serverless local HLLCounters and Sets")
	fmt.Println("Use -b </path/to/databasename.db> to persist w/ badgedb (default name nazare.db in the current dir)")
	fmt.Println("HyperLogLog based counters:")
	fmt.Println("Add to a hll counter: nazare-cli counters -a <countername> <item>")
	fmt.Println("Estimate counter size: nazare-cli -c -e <countername>")
	fmt.Println("Cuckoo filter based sets:")
	fmt.Println("Add to a set: nazare-cli -s -a <setname> <item>")
	fmt.Println("Check if an item belongs to a set: nazare-cli -s -i <setname> <item>")
	fmt.Println("Remove an item from a set: nazare-cli -s -r <setname> <item>")
	fmt.Println("Estimate set cardinality: nazare-cli -s -c <setname>")
	fmt.Println("That's it, no way to get an item from a set, cuckoo filter stores hashes and signal it an item was 'seen'")
	fmt.Println("K/V handling:")
	fmt.Println("Set a key: nazare-cli -k -s <keyname> <value>")
	fmt.Println("Get the value of a key: nazare-cli -k -g <keyname>")
	fmt.Println("Delete a key: nazare-cli -k -d <keyname>")

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
		badgerds, _ := db.NewBadgerDatastorage("testcounters.db")
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
