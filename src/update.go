package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func CheckForPeekUpdate() {
	var latestVersion string
	var currentVersion string

	currentVersion = VERSION

	log.Info("Update checker initialized!")

	for {
		// log
		log.Info("Checking for updates...")

		// make a GET request to GitHub api to get all releases
		url := "https://api.github.com/repos/fwuffyboi/peek/releases"
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("Error making request to GitHub API: %s", err)
		}
		defer resp.Body.Close() // todo error handle this

		// read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading response body: %s", err)
		}

		// parse response
		var releases []Release
		if err := json.Unmarshal(body, &releases); err != nil {

			// if error, see if it's a ratelimit error.
			if strings.Contains(string(body), "API rate limit exceeded") {
				log.Warnf("GitHub API rate limit exceeded! Will check for updates again in 1 hour...")

				// send alert
				if _, ok := getAlerts()["GitHub API rate limit exceeded!"]; !ok {
					addAlert("GitHub API rate limit exceeded!")
				} else {
					log.Info("Alert already sent, skipping...")
				}
			} else {
				log.Errorf("Unknown error unmarshalling JSON: %s", err)
			}

		} else {
			// Unmarshalled it just fine. Continue with the rest of the code.

			if len(releases) == 0 {
				log.Info("No releases found! Assuming on latest version.")
				log.Info("Will check for updates in 1 hour...")
			} else {
				// compare versions
				latestVersion = releases[0].TagName
				relComp := compareVersions(currentVersion, latestVersion)
				if relComp == -1 {
					// current version is older than newer version
					log.Infof("New version found! Current: %s, Latest: %s", currentVersion, latestVersion)

					// check if alert already sent
					if _, ok := getAlerts()["New update available! Current: "+currentVersion+", Latest: "+latestVersion]; !ok {
						addAlert("New update available! Current: " + currentVersion + ", Latest: " + latestVersion)
					} else {
						log.Info("Alert already sent, skipping...")
					}

				} else if relComp == 1 {
					log.Warnf("You are running a newer version than the latest release! Current: %s, Latest: %s", currentVersion, latestVersion)

					// check if alert already sent
					if _, ok := getAlerts()["You are running a newer version than the latest release! Current: "+currentVersion+", Latest: "+latestVersion]; !ok {
						addAlert("You are running a newer version than the latest release! Current: " + currentVersion + ", Latest: " + latestVersion)
					} else {
						log.Info("Alert already sent, skipping...")
					}
				} else {
					log.Info("You are running the latest version!")
				}

				log.Info("Will check for updates again in 1 hour...")
			}
		}

		time.Sleep(1 * time.Hour)
	}
}

// Chatgpt shit
func compareVersions(currentVersion, latestVersion string) int {
	currentParts := strings.Split(currentVersion, ".")
	latestParts := strings.Split(latestVersion, ".")

	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		currentPart := currentParts[i]
		latestPart := latestParts[i]

		if currentPart == latestPart {
			continue
		}

		currentNumber := 0
		latestNumber := 0

		fmt.Sscanf(currentPart, "%d", &currentNumber) // todo error handle
		fmt.Sscanf(latestPart, "%d", &latestNumber)   // todo error handle

		if currentNumber < latestNumber {
			return -1 // currentVersion is older
		} else if currentNumber > latestNumber {
			return 1 // currentVersion is newer
		}
	}

	if len(currentParts) < len(latestParts) {
		return -1 // currentVersion is older
	} else if len(currentParts) > len(latestParts) {
		return 1 // currentVersion is newer
	}

	return 0 // versions are the same
}
