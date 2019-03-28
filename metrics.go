package main

import (
	gometrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
	"time"
)

type LocalMetrics struct {
	MessageCount gometrics.Counter
	ClientCount  gometrics.Counter
}

var metrics = &LocalMetrics{
	MessageCount: gometrics.NewCounter(),
	ClientCount:  gometrics.NewCounter(),
}

func (m *LocalMetrics) Report() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.WithField("MessageCount", m.MessageCount.Count()).WithField("ClientCount", m.ClientCount.Count()).Info("report")
		}
	}
}
