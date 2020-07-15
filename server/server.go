package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	httplogger "github.com/gleicon/go-httplogger"
	"github.com/gleicon/nazare/datalayer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tidwall/redcon"
)

/*
NZServer contains all the server needs to run
*/
type NZServer struct {
	ldb *datalayer.LocalDB

	serverAddr  string
	httpAPIAddr string
	dbPath      string

	customMetrics *NZServerCustomMetrics
}

/*
NewNZServer starts a new nazare server (redis and http endpoints)
*/
func NewNZServer(serverAddr, httpAPIAddr, dbPath string) (*NZServer, error) {
	var err error
	if dbPath == "" {
		return nil, errors.New("No database path given")
	}
	nzServer := NZServer{serverAddr: serverAddr, httpAPIAddr: httpAPIAddr, dbPath: dbPath}
	nzServer.ldb = datalayer.NewLocalDB()

	log.Println("Creating local database")
	nzServer.ldb.Start(dbPath)

	log.Println("Creating metrics")
	if nzServer.customMetrics, err = NewNZServerCustomMetrics(); err != nil {
		return nil, errors.Unwrap(fmt.Errorf("Error creating metrics pool: %w", err))
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
