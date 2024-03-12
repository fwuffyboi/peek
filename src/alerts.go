package main

import "time"

func addAlert(message string) {
	alertsList[message] = time.Now()
}

func getAlerts() map[string]time.Time {
	return alertsList
}
