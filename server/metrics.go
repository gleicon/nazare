package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/*
NZServerCustomMetrics is a local metrics pool
*/
type NZServerCustomMetrics struct {
	successfulCmds    prometheus.Counter
	errorCmds         prometheus.Counter
	totalClients      prometheus.Counter
	totalRequests     prometheus.Counter
	activeClients     prometheus.Gauge
	uninmplementedCmd prometheus.Counter
}

/*
NewNZServerCustomMetrics create a new metric pool
*/
func NewNZServerCustomMetrics() (*NZServerCustomMetrics, error) {
	nm := NZServerCustomMetrics{}

	nm.successfulCmds = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nazare_succesfully_processed_commands",
		Help: "The total number of processed commands without error",
	})

	nm.errorCmds = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nazare_commands_with_error",
		Help: "The total number of processed commands that yielded an error",
	})

	nm.totalClients = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nazare_total_clients_ever",
		Help: "The total number of clients that connected to this instance",
	})

	nm.totalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nazare_total_requests",
		Help: "The total number of requests processed by this instance",
	})
	nm.activeClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nazare_active_clients",
		Help: "The number of active clients connected",
	})
	nm.uninmplementedCmd = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nazare_uninmplemented_command",
		Help: "The total calls to commands not implemented",
	})
	return &nm, nil
}
