package counters

import (
	"errors"
	"sync"

	"github.com/axiomhq/hyperloglog"
	"github.com/gleicon/nazare/db"
)

/*
HLLCounterStats general stats about a set of counters and its db
*/
type HLLCounterStats struct {
	ActiveCounters uint64
}

/*
HLLCounters Wraps hyperloglog implementation and keep an active counter register
*/
type HLLCounters struct {
	hllrwlocks  map[string]*sync.RWMutex
	datastorage db.Datastorage
	stats       HLLCounterStats
}

/*
NewHLLCounters is a high level HLL counter abstraction.
Wraps the hll implementation, keeps an active counter register, knows how to save it to disk.
*/
func NewHLLCounters(ds db.Datastorage) (*HLLCounters, error) {

	hll := HLLCounters{}
	hll.datastorage = ds
	hll.hllrwlocks = make(map[string]*sync.RWMutex)

	return &hll, nil
}

/*
Close flushes and closed the database
*/
func (hc *HLLCounters) Close() {
	// TODO: this is the place to flush what we have for the db.
	hc.datastorage.Close()
}

/*
Stats for all counters
*/
func (hc *HLLCounters) Stats() HLLCounterStats {
	// TODO: this is the place to flush what we have for the db.
	return hc.stats
}

/*
RetrieveCounterEstimate retrieves the estimate for <<name>> counter
*/
func (hc *HLLCounters) RetrieveCounterEstimate(name string) (uint64, error) {
	var err error
	var cc []byte
	if cc, err = hc.datastorage.Get(name); err != nil {
		return 0, err
	}
	if cc == nil {
		return 0, errors.New("Counter does not exist:" + name)
	}

	hll := hyperloglog.New16()
	if err := hll.UnmarshalBinary(cc); err != nil {
		return 0, err
	}
	return hll.Estimate(), nil
}

/*
RetrieveAndMergeCounterEstimates retrieves the estimate for <<name>> counter
*/
func (hc *HLLCounters) RetrieveAndMergeCounterEstimates(counterNames ...string) (uint64, error) {
	var err error
	var cc []byte
	pivotHLL := hyperloglog.New16()
	for _, counter := range counterNames {
		if cc, err = hc.datastorage.Get(counter); err != nil {
			return 0, err
		}

		if cc == nil {
			continue // just skip or
			// return 0, errors.New("Counter does not exist:" + name)
		}
		tempHLL := hyperloglog.New16()
		if err := tempHLL.UnmarshalBinary(cc); err != nil {
			return 0, err
		}
		if err := pivotHLL.Merge(tempHLL); err != nil {
			return 0, nil
		}
	}
	return pivotHLL.Estimate(), nil
}

/*
IncrementCounter increments <<name>> counter by adding <<item>> to it.
The counter and its lock is automatically created if it is empty.
*/
func (hc *HLLCounters) IncrementCounter(name string, item []byte) error {
	if hc.hllrwlocks[name] == nil {
		hc.hllrwlocks[name] = new(sync.RWMutex)
	}

	localMutex := hc.hllrwlocks[name]
	localMutex.Lock()
	defer localMutex.Unlock()

	cc, _ := hc.datastorage.Get(name)

	hll := hyperloglog.New16()
	if cc != nil {
		if err := hll.UnmarshalBinary(cc); err != nil {
			return err
		}
	} else {
		hc.stats.ActiveCounters++
	}

	hll.Insert(item)
	var bd []byte
	var err error

	if bd, err = hll.MarshalBinary(); err != nil {
		return err
	}

	if err = hc.datastorage.Add(name, bd); err != nil {
		return err
	}
	return nil
}

/*
ExportCounter returns the counter binary form
*/
func (hc *HLLCounters) ExportCounter(name string) ([]byte, error) {
	cc, _ := hc.datastorage.Get(name)
	if cc == nil {
		return nil, errors.New("Counter does not exist:" + name)
	}
	hll := hyperloglog.New16()
	if err := hll.UnmarshalBinary(cc); err != nil {
		return nil, err
	}
	return hll.MarshalBinary()
}
