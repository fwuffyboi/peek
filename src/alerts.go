package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func addAlert(message string) {
	var timeRN = time.Now()
	alertsList[message] = timeRN
	log.Infof("Added alert: message: %s, time: %q", message, timeRN)

	// Send message to Telegram user
	config, err := ConfigParser()
	if err != nil {
		log.Error("Error reading config file: ", err)
		return
	}

	if !config.Integrations.Telegram.Enabled {
		log.Warn("Telegram integration not enabled. Will not attempt to send message.")
		return
	} else {

		// the below line replaces the %alert% placeholder in the config file with the actual
		// alert message, so users can customise the message sent to themselves.
		var tgMessage = strings.Replace(config.Integrations.Telegram.TelegramMessage, "%alert%", message, -1)

		err = sendTelegramMessage(tgMessage)
		if err != nil {
			log.Error("Error sending message to Telegram: ", err)
		}
	}
}

func getAlerts() map[string]time.Time {
	return alertsList
}

// @Summary Returns all alerts
// @Description Returns all alerts
// @Produce json
// @Success 200 {object} map[string]time.Time
// @Failure 429
// @Failure 500 {string} string "Internal Server Error"
// @Tags apiPeekGroup
// @Router /peek/alerts [get]
func apiReturnAlerts(c *gin.Context) {

	// todo: add auth check

	c.JSON(200, getAlerts())
}
