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
	return hc.stats
}

/*
RetrieveCounterEstimate retrieves the estimate for <<name>> counter
*/
func (hc *HLLCounters) RetrieveCounterEstimate(key []byte) (uint64, error) {
	var err error
	var cc []byte
	if cc, err = hc.datastorage.Get(key); err != nil {
		return 0, err
	}
	if cc == nil {
		return 0, errors.New("Counter does not exist:" + string(key))
	}

	hll := hyperloglog.New16()
	if err := hll.UnmarshalBinary(cc); err != nil {
		return 0, err
	}
	return hll.Estimate(), nil
}

/*
hll aware merge function. problem, it ties the db impl w/ the hll impl
*/
func (hc *HLLCounters) mergeMarshaledCounter(dest []byte, increment []byte) ([]byte, error) {
	hll := hyperloglog.New16()
	if hll != nil {
		if err := hll.UnmarshalBinary(dest); err != nil {
			return nil, err
		}
	}
	hll.Insert(increment)
	return hll.MarshalBinary()
}

/*
RetrieveAndMergeCounterEstimates retrieves the estimate for <<name>> counter
*/
func (hc *HLLCounters) RetrieveAndMergeCounterEstimates(counterNames ...[]byte) (uint64, error) {
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
The naive implementation locks(), get, increment and set
The counter and its lock are automatically created if it is empty.
*/
func (hc *HLLCounters) IncrementCounter(key []byte, item []byte) error {
	if hc.hllrwlocks[string(key)] == nil {
		hc.hllrwlocks[string(key)] = new(sync.RWMutex)
	}

	localMutex := hc.hllrwlocks[string(key)]
	localMutex.Lock()
	defer localMutex.Unlock()

	cc, _ := hc.datastorage.Get([]byte(key))

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

	if err = hc.datastorage.Add([]byte(key), bd); err != nil {
		return err
	}
	return nil
}

/*
ExportCounter returns the counter binary form
*/
func (hc *HLLCounters) ExportCounter(key string) ([]byte, error) {
	cc, _ := hc.datastorage.Get([]byte(key))
	if cc == nil {
		return nil, errors.New("Counter does not exist:" + key)
	}
	hll := hyperloglog.New16()
	if err := hll.UnmarshalBinary(cc); err != nil {
		return nil, err
	}
	return hll.MarshalBinary()
}
