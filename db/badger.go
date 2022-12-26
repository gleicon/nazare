package db

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/golang/snappy"
)

/*
BadgerDatastorage is an abstraction over badger db to store bytemaps of serialised hll counters
*/
type BadgerDatastorage struct {
	Datastorage
	filepath    string
	datapath    string
	db          *badger.DB
	compression bool
}

/*
NewBadgerDatastorage spin up a new badge based datastorage
*/
func NewBadgerDatastorage(dbpath string) (*BadgerDatastorage, error) {
	var err error
	bds := BadgerDatastorage{filepath: dbpath, datapath: dbpath}
	opts := badger.DefaultOptions(dbpath)
	opts.SyncWrites = true
	opts.Dir = dbpath
	bds.db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			log.Println("Database housekeeping")
			bds.cleanup()
			time.Sleep(5 * time.Minute)
		}
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
	bds.db.RunValueLogGC(1.0)
}

/*
Add data
*/
func (bds *BadgerDatastorage) Add(key, value []byte) error {
	var vals []byte
	if bds.compression {
		vals := snappy.Encode(nil, value)

		if vals == nil {
			return errors.New("Error compressing payload")
		}
	} else {
		vals = value
	}

	err := bds.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), vals)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

/*
Merge value into key, this require knowledge about the hll function and no compression
*/
func (bds *BadgerDatastorage) Merge(key, value []byte, mergeFunc badger.MergeFunc) error {
	if bds.compression {
		return errors.New("Merge requires non-compressed payload")
	}

	// This may be a thing to tune. compaction goroutine from badger is set to 200ms
	dm := bds.db.GetMergeOperator(key, mergeFunc, 200*time.Millisecond)
	defer dm.Stop()

	dm.Add(value)

	err := bds.db.Update(func(txn *badger.Txn) error {
		lv, _ := dm.Get()
		err := txn.Set([]byte(key), lv)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

/*
Get data from the db
*/
func (bds *BadgerDatastorage) Get(key []byte) ([]byte, error) {
	var payload []byte
	err := bds.db.View(func(txn *badger.Txn) error {

		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		err = item.Value(func(val []byte) error {
			if !bds.compression {
				payload = val
			} else {
				retVal, err := snappy.Decode(nil, val)
				if err != nil {
					return err
				}
				if payload == nil {
					return errors.New("Error uncompressing payload")
				}
				payload = retVal
			}
			return nil

		})
		return err
	})

	return payload, err
}

/*
Delete a given key
*/
func (bds *BadgerDatastorage) Delete(key []byte) (bool, error) {
	var found bool
	err := bds.db.Update(func(txn *badger.Txn) error {

		err := txn.Delete(key)
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			found = false
			return nil
		}
		found = true
		return nil
	})
	return found, err
}

/*
Close the database safely, flushing it first
*/
func (bds *BadgerDatastorage) Close() {
	bds.cleanup()
	bds.db.Close()
}

/*
Flush the database safely
*/
func (bds *BadgerDatastorage) Flush() {
	bds.cleanup()
}
