package main

import (
	"log"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
)

/*
BadgerDatastorage is an abstraction over badger db to store bytemaps of serialised hll counters
*/
type BadgerDatastorage struct {
	Datastorage
	filepath string
	datapath string
	db       *badger.DB
}

/*
NewBadgerDatastorage spin up a new badge based datastorage
*/
func NewBadgerDatastorage(dbpath string) (*BadgerDatastorage, error) {
	var err error
	bds := BadgerDatastorage{filepath: dbpath, datapath: dbpath}
	opts := badger.DefaultOptions
	opts.SyncWrites = true
	opts.Dir = dbpath
	opts.ValueDir = dbpath
	bds.db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}
	// read the counter database on load. That could (should ?) be moved
	//if err = bds.loadCountersFromDatabase(); err != nil {
	//	return nil, err
	//}
	go func() {
		log.Println("Database housekeeping")
		bds.cleanup()
		time.Sleep(5 * time.Minute)
	}()
	return &bds, nil
}

/*
Cleanup keeps the db healthy
*/
func (bds *BadgerDatastorage) cleanup() {
	var ll sync.Mutex
	ll.Lock()
	defer ll.Unlock()
	// TODO: these cleanup methods should be in a goroutine
	bds.db.PurgeOlderVersions()
	bds.db.RunValueLogGC(1.0)
}

func (bds *BadgerDatastorage) Add(key string, value []byte) error {
	err := bds.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), value)
		if err != nil {
			return err
		}
		return nil
	})
	return err

}

func (bds *BadgerDatastorage) Get(key string) ([]byte, error) {
	var payload []byte
	err := bds.db.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		payload, _ = item.Value()
		return nil
	})

	return payload, err
}

func (bds *BadgerDatastorage) Close() {
	bds.cleanup()
	bds.db.Close()
}

func (bds *BadgerDatastorage) Flush() {
	bds.cleanup()
}

func Delete(string) {}

/*
Read all counters from the database, overwrites any current value from memory
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
*/
