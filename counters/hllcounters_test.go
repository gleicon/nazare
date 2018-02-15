package counters

import (
	"fmt"
	"os"
	"testing"

	"github.com/gleicon/nazare/db"
)

var counters *HLLCounters

func TestMain(t *testing.T) {
	memds, _ := db.NewHLLDatastorage()
	counters, _ = NewHLLCounters(memds)
}
func TestRetrieveCounterEstimate(t *testing.T) {
	var cc uint64
	var err error
	counterName := "testcounter"

	if cc, err = counters.RetrieveCounterEstimate(counterName); err == nil {
		t.Error(err)
		return
	}

	if err := counters.IncrementCounter(counterName, []byte("testIncrement")); err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	if cc, err = counters.RetrieveCounterEstimate(counterName); err != nil {
		t.Error(err)
		return
	}

	if cc != 1 {
		t.Error("Wrong estimate")
		return
	}

}
