package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func addAlert(message string) {
	var timeRN = time.Now()
	alertsList[message] = timeRN
	log.Infof("Added alert: message: %s, time: %q", message, timeRN)
}

func getAlerts() map[string]time.Time {
	return alertsList
}

func apiReturnAlerts(c *gin.Context) {
	c.JSON(200, getAlerts())
}
