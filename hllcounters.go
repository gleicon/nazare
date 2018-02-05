package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/dgraph-io/badger"
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
	hllcounters map[string]*hyperloglog.Sketch
	hllrwlocks  map[string]*sync.RWMutex
	db          *badger.DB
	stats       HLLCounterStats
}

/*
NewHLLCounters is a high level HLL counter abstraction.
Wraps the hll implementation, keeps an active counter register, knows how to save it to disk.
*/
func NewHLLCounters(dbpath string) (*HLLCounters, error) {
	var err error

	hll := HLLCounters{}
	hll.hllcounters = make(map[string]*hyperloglog.Sketch)
	hll.hllrwlocks = make(map[string]*sync.RWMutex)

	opts := badger.DefaultOptions
	opts.SyncWrites = true
	opts.Dir = dbpath
	opts.ValueDir = dbpath
	hll.db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}
	// read the counter database on load. That could (should ?) be moved
	if err = hll.loadCountersFromDatabase(); err != nil {
		return nil, err
	}
	go func() {
		log.Println("Cleaning up database")
		hll.cleanup()
		time.Sleep(5 * Duration.Minute)
	}()
	return &hll, nil
}

/*
Close flushes and closed the database
*/
func (hc *HLLCounters) Close() {
	// TODO: this is the place to flush what we have for the db.
	hc.cleanup()
	hc.db.Close()
}

/*
Cleanup keeps the db healthy
*/
func (hc *HLLCounters) cleanup() {
	ll := sync.Mutex
	ll.Lock()
	defer ll.Unlock()
	// TODO: these cleanup methods should be in a goroutine
	hc.db.PurgeOlderVersions()
	hc.db.RunValueLogGC()
}

/*
Stats for all counters
*/
func (hc *HLLCounters) Stats() HLLCounterStats {
	// TODO: this is the place to flush what we have for the db.
	return hc.stats
}

/*
Read all counters from the database, overwrites any current value from memory
*/
func (hc *HLLCounters) loadCountersFromDatabase() error {
	log.Println("Loading counters from database")
	hc.stats.ActiveCounters = 0
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err != nil {
				return err
			}
			key := string(k)
			hl := hyperloglog.New16()
			hl.UnmarshalBinary(v)
			hc.hllcounters[key] = &hl
			hc.stats.ActiveCounters++
		}
		return nil
	})
	return err
}

/*
RetrieveCounterEstimate retrieves the estimate for <<name>> counter
*/
func (hc *HLLCounters) RetrieveCounterEstimate(name string) (uint64, error) {
	cc := hc.hllcounters[name]
	if cc == nil {
		return 0, errors.New("Counter does not exist:" + name)
	}
	return cc.Estimate(), nil
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

	cc := hc.hllcounters[name]
	if cc == nil {
		hc.hllcounters[name] = hyperloglog.New16()
		hc.stats.ActiveCounters++
	}
	hc.hllcounters[name].Insert(item)
	err := hc.db.Update(func(txn *badger.Txn) error {
		item, err = txn.Get([]byte(name))
		// TODO: item exists ? diff from memory ?
		//TODO set new value
		return nil
	})
	return err
}

/*
ExportCounter returns the counter binary form
*/
func (hc *HLLCounters) ExportCounter(name string) ([]byte, error) {
	cc := hc.hllcounters[name]
	if cc == nil {
		return nil, errors.New("Counter does not exist: " + name)
	}
	return cc.MarshalBinary()
}
