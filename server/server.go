package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	httplogger "github.com/gleicon/go-httplogger"
	"github.com/gleicon/nazare/counters"
	"github.com/gleicon/nazare/db"
	"github.com/gleicon/nazare/sets"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tidwall/redcon"
)

/*
NZServer contains all the server needs to run
*/
type NZServer struct {
	localCounters    *counters.HLLCounters
	localSets        *sets.CkSet
	localDatastorage db.Datastorage

	serverAddr  string
	httpAPIAddr string
	dbPath      string
}

/*
NewNZServer starts a new nazare server (redis and http endpoints)
*/
func NewNZServer(serverAddr, httpAPIAddr, dbPath string) (*NZServer, error) {
	var err error
	nzServer := NZServer{serverAddr: serverAddr, httpAPIAddr: httpAPIAddr, dbPath: dbPath}

	log.Println("Creating local database")
	if nzServer.localDatastorage, err = db.NewBadgerDatastorage(dbPath); err != nil {
		return nil, errors.Unwrap(fmt.Errorf("Error creating datastorage: %w", err))
	}

	log.Println("Creating counter structures")
	if nzServer.localCounters, err = counters.NewHLLCounters(nzServer.localDatastorage); err != nil {
		return nil, errors.Unwrap(fmt.Errorf("Error creating HLLCounters: %w", err))

	}

	log.Println("Creating sets structures")
	if nzServer.localSets, err = sets.NewCkSets(nzServer.localDatastorage); err != nil {
		return nil, errors.Unwrap(fmt.Errorf("Error creating localDatastorage: %w", err))

	}
	return &nzServer, nil
}

/*
Start the service
*/
func (nzs *NZServer) Start() error {
	var err error
	// spin up the Redis connector
	errChannel := make(chan error, 1)
	go func() {
		defer close(errChannel)
		log.Println("Redis protocol server: " + nzs.serverAddr)

		err := redcon.ListenAndServe(nzs.serverAddr, nzs.redisCommandParser, nzs.newConnection, nzs.closeConnection)
		if err != nil {
			log.Println("Error spinning up Redis protocol connector: ", err)
			errChannel <- errors.Unwrap(fmt.Errorf("Error spinning up Redis connector: %w", err))
		}
		select {}
	}()

	log.Println("HTTP Metrics API: " + nzs.httpAPIAddr)

	if err := nzs.httpMetrics(); err != nil {
		return err
	}
	err = <-errChannel
	return err
}

func (nzs *NZServer) httpMetrics() error {
	serveMux := http.NewServeMux()

	serveMux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    nzs.httpAPIAddr,
		Handler: httplogger.HTTPLogger(serveMux),
	}
	if err := srv.ListenAndServe(); err != nil {
		return errors.Unwrap(fmt.Errorf("Error spinning up HTTP Metrics API: %w", err))
	}
	return nil
}
