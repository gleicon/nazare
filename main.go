package main

import (
	"flag"
	"fmt"
	"os"
)

func help() {
	fmt.Println("increment hll counters: counters -n <countername> -i foobar")
	fmt.Println("estimate counter size: counters -n <countername>")
	os.Exit(-1)
}

func main() {

	memds, _ := NewHLLDatastorage()
	counters, _ := NewHLLCounters(memds)
	counterName := flag.String("n", "", "HLL counter name")
	incrementBy := flag.String("i", "", "Increment counter by string")
	flag.Usage = help
	flag.Parse()

	if *counterName == "" {
		fmt.Println("No counter")
		os.Exit(-1)
	}

	if *incrementBy == "" {
		var cc uint64
		var err error
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
		var cc uint64
		var err error
		if cc, err = counters.RetrieveCounterEstimate(*counterName); err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)

		}
		fmt.Printf("counter: %s - est. size: %d\n", *counterName, cc)
		xx, _ := counters.ExportCounter(*counterName)
		fmt.Printf("%x\n", xx)

	}

}
