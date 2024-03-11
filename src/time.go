package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func serverTimezone() string {
	t := time.Now()
	zone, offset := t.Zone()
	log.Infof("Current server time: %s. Current server time offset: %d", zone, offset)

	return zone
}
