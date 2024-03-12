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

		// parse it
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading response body: %s", err)
		}

		var releases []Release
		if err := json.Unmarshal(body, &releases); err != nil {
			log.Errorf("Error unmarshalling JSON: %s", err)
		}

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
				if _, ok := getAlerts()["New version found! Current: "+currentVersion+", Latest: "+latestVersion]; !ok {
					addAlert("New update available! Current: " + currentVersion + ", Latest: " + latestVersion)
				} else {
					log.Info("Alert already sent, skipping...")
				}

			} else if relComp == 1 {
				log.Infof("You are running a newer version than the latest release! How are you even doing this! Current: %s, Latest: %s", currentVersion, latestVersion)
				addAlert("You are running a newer version than the latest release! How are you even doing this! Current: " + currentVersion + ", Latest: " + latestVersion)
			} else {
				log.Info("You are running the latest version!")
			}

			// else, sleep for 1 hour
			log.Infof("No updates found, sleeping for 1 hour...")
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
