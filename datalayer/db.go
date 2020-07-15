package datalayer

import (
	"errors"
	"fmt"
	"log"

	"github.com/gleicon/nazare/counters"
	"github.com/gleicon/nazare/db"
	"github.com/gleicon/nazare/sets"
)

// LocalDB is the minimal set of resources for nazare. Works for both cli and server
type LocalDB struct {
	dbPath           string
	LocalDatastorage db.Datastorage
	LocalCounters    *counters.HLLCounters
	LocalSets        *sets.CkSet
}

// NewLocalDB preps the local db
func NewLocalDB() *LocalDB {
	ldb := LocalDB{}
	return &ldb
}

// Start spins a local structure
func (ldb *LocalDB) Start(dbPath string) error {
	var err error

	if ldb.LocalDatastorage, err = db.NewBadgerDatastorage(dbPath); err != nil {
		return errors.Unwrap(fmt.Errorf("Error creating datastorage: %w", err))
	}

	log.Println("Creating counter structures")
	if ldb.LocalCounters, err = counters.NewHLLCounters(ldb.LocalDatastorage); err != nil {
		return errors.Unwrap(fmt.Errorf("Error creating HLLCounters: %w", err))

	}

	log.Println("Creating sets structures")
	if ldb.LocalSets, err = sets.NewCkSets(ldb.LocalDatastorage); err != nil {
		return errors.Unwrap(fmt.Errorf("Error creating localDatastorage: %w", err))

	}

	return nil
}
