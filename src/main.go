package main

import (
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

// Define constants
const (
	// DefaultWebuiAddress Default address for the web UI, in case it is not provided/invalid in the config file
	DefaultWebuiAddress = "0.0.0.0:42649" // Address of the webserver, HAS to be in the format of: IP:PORT

	// VERSION Version of Peek
	VERSION = "v0.9.7-alpha" // Version of Peek

	// DefaultWebUiHost Default host for the web UI, in case it is not provided/invalid in the config file
	DefaultWebUiHost = "0.0.0.0"
	// DefaultWebUiPort Default port for the web UI, in case it is not provided/invalid in the config file
	DefaultWebUiPort = 42649
)

var ServerIPAddress = "" // IP address of the server
var ServerCountry = ""   // Country of the server, based on IP

var alertsList = make(map[string]time.Time) // List of alerts

var authTokens []string // List of auth tokens

func main() {
	// Setup logging and obtain the log file handle and multi-writer
	_, _, err := setupLogging()
	if err != nil {
		log.Fatalf("Failed to setup logging: %s", err)
		return
	}

	// Log startup info
	startupLogs()

	// check if OS is supported (linux)
	isOSSupported()

	// start the update checker thread
	go CheckForPeekUpdate()

	// Get the server ip and save into var
	log.Info("Attempting to get server's IP address.")
	ServerIPAddress = getIP()

	log.Infof("Attempting to get server's country from IP address.")
	ServerCountry = countryFromIP(ServerIPAddress)

	// Get IP and port to run webserver on
	config, err := ConfigParser()
	if err != nil {
		log.Fatalf("Failed to get config: %s", err)
	}
	host := config.Api.ApiHost
	port := config.Api.ApiPort
	ginRatelimit := config.Api.RateLimit

	// Run the webserver
	runGin(host, port, ginRatelimit)
}

func startupLogs() {
	// Log when shit was started up
	log.Infof("Peek Version: %s", VERSION)
	log.Infof("Application started at time: %s, on the date: %s(YYYY-MM-DD).",
		time.Now().Format("15:04:05"),   // Format for hh:mm:ss
		time.Now().Format("2006-01-02"), // Format for yyyy-mm-dd
	)

	// Copyright notices
	log.Info("Peek is licensed under the MIT License. See LICENSE for more information.")
	log.Infof("(C) %s Peek Contributors. All rights reserved.", time.Now().Format("2006"))
	log.Info("(C) IP Geolocation by DB-IP, https://db-ip.com/")
	log.Info("Peek is a program written by: @fwuffyboi (https://github.com/fwuffyboi)")

}

func isOSSupported() bool {
	// Check os, if windows/darwin, panic, else continue
	log.Info("Checking what operating system is in use...")
	switch runtime.GOOS {
	case "linux":
		log.Info("Linux derivative detected.")
		return true
	default:
		log.Error("Unsupported operating system detected.")
		log.Fatal("This program only supports Linux distributions.")
		return false
	}
}
