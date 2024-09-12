package vip

import (
	"github.com/jotak/net-infra-mon/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	vipHostsDef = metrics.Declare(
		"vip_hosts",
		"Number of virtual IPs per host",
		metrics.TypeGauge,
		"vip",
		"host",
	)
	vipEventsDef = metrics.Declare(
		"vip_events",
		"Number and nature of virtual IP events per host",
		metrics.TypeCounter,
		"vip",
		"host",
		"event",
	)
)

///////////////////////////////////////////////////////////////////////////////
// Helper structs

type eventCounter struct{ *prometheus.CounterVec }

func createEventsCounter() eventCounter {
	return eventCounter{metrics.NewCounterVec(&vipEventsDef)}
}

func (c *eventCounter) increase(vip, host, event string) {
	c.WithLabelValues(vip, host, event).Add(1)
}

type hostsGauge struct{ *prometheus.GaugeVec }

func createVIPHostsGauge() hostsGauge {
	return hostsGauge{metrics.NewGaugeVec(&vipHostsDef)}
}

func (g *hostsGauge) set(vip, host string) {
	g.WithLabelValues(vip, host).Set(1)
}
