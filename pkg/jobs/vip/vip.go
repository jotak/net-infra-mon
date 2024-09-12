package vip

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "vip")

func Run(_ context.Context) {
	host := getHostName()
	vipEventsCounter := createEventsCounter()
	vipHostsGauge := createVIPHostsGauge()

	// TODO: list VIPs
	for _, vip := range []string{"foo", "bar"} {
		vipHostsGauge.set(vip, host)
		// TODO: listen to event
		vipEventsCounter.increase(vip, host, "an-event")
	}
}

func getHostName() string {
	// TODO: should it be inferred from config?
	hostname, err := os.Hostname()
	if err != nil {
		log.WithError(err).Error("could not get host name")
		return "unknown"
	}
	return hostname
}
