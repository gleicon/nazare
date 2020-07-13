package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/gleicon/nazare/counters"
	"github.com/gleicon/nazare/db"
	"github.com/gleicon/nazare/sets"
)

// LocalDb is the minimal set of resources for nazare
type LocalDb struct {
	dbPath           string
	localDatastorage db.Datastorage
	localCounters    *counters.HLLCounters
	localSets        *sets.CkSet
}

// NewLocalDB preps the local db
func NewLocalDB() *LocalDb {
	ldb := LocalDb{}
	return &ldb
}

// Start spins a local structure
func (ldb *LocalDb) Start(dbPath string) error {
	var err error

	if ldb.localDatastorage, err = db.NewBadgerDatastorage(dbPath); err != nil {
		return errors.Unwrap(fmt.Errorf("Error creating datastorage: %w", err))
	}

	log.Println("Creating counter structures")
	if ldb.localCounters, err = counters.NewHLLCounters(ldb.localDatastorage); err != nil {
		return errors.Unwrap(fmt.Errorf("Error creating HLLCounters: %w", err))

	}

	log.Println("Creating sets structures")
	if ldb.localSets, err = sets.NewCkSets(ldb.localDatastorage); err != nil {
		return errors.Unwrap(fmt.Errorf("Error creating localDatastorage: %w", err))

	}

	return nil
}
