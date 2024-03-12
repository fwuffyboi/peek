package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func addAlert(message string) {
	alertsList[message] = time.Now()
}

func getAlerts() map[string]time.Time {
	return alertsList
}

func apiReturnAlerts(c *gin.Context) {
	c.JSON(200, getAlerts())
}
