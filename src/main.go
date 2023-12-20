package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"time"
)

// Define constants
const (
	DefaultWebUiAddr = "0.0.0.0:42649" // Address of the webserver, HAS to be in the format of: IP:PORT
	VERSION          = "0.0.1"         // Version of Peek
	DefaultWebUiHost = "0.0.0.0"
	DefaultWebUiPort = 42649
)

var UnsupportedOS = false // assume false until proven true
var IpAddress = ""        // IP address of the server
var ServerCountry = ""    // Country of the server, based on IP

func main() {
	// Log when shit was started up
	log.Infof("Peek Version: %s", VERSION)
	log.Infof("Application started at time: %s, on the date: %s(YYYY-MM-DD).",
		time.Now().Format("15:04:05"),   // Format for hh:mm:ss
		time.Now().Format("2006-01-02"), // Format for yyyy-mm-dd
	)

	// Copyright notices
	log.Info("Peek is licensed under the MIT License. See LICENSE for more information.")
	log.Infof("(C) %s Peek Contributors. All rights reserved.", time.Now().Format("2006"))
	log.Info("Peek is a program written by: @fwuffyboi (https://github.com/fwuffyboi)")
	log.Info("") // free space
	log.Info("(C) IP Geolocation by DB-IP, https://db-ip.com/")
	log.Info("") // free space

	// Check os, if windows, panic, else continue
	log.Info("Checking what operating system is in use...")
	switch runtime.GOOS {
	case "linux", "darwin":
		log.Info("Linux/darwin derivative detected.")
	default:
		log.Error("Unsupported operating system detected.")
		log.Panic("This program only supports the Linux/Darwin(MacOS) operating systems.")
		UnsupportedOS = true
	}

	// Get the server ip and save into var
	IpAddress = getIP()

	log.Infof("Attempting to get server's country from IP address.")
	ServerCountry = countryFromIP(IpAddress)

	// Get IP and port to run webserver on
	config, err := ConfigParser()
	if err != nil {
		log.Fatalf("Failed to get config: %s", err)
	}
	host := config.Api.ApiHost
	port := config.Api.ApiPort

	// Run the webserver
	runGin(host, port)
}

// Starts and runs the webserver, using the gin framework
func runGin(host string, port int) {
	log.Info("GIN: Starting the Peek WebServer...")

	gin.SetMode(gin.ReleaseMode) // set to production mode
	r := gin.Default()
	r.ForwardedByClientIP = true

	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Panicf("GIN: Failed to set trusted proxies: %s", err)
	}

	// Define routes
	// INFO routes
	r.GET("/api/", func(c *gin.Context) { apiEndpoints(c) }) // return all api endpoints

	// NOACTION routes
	r.GET("/api/full/", func(c *gin.Context) { apiFull(c) })             // Return all api/json info
	r.GET("/api/disk/all/", func(c *gin.Context) { storageAllDisks(c) }) // Return all disk info
	// AFTER TESTING OUTSIDE OF FLATPAK, THE ABOVE ROUTE WORKS. BUT ONLY OUTSIDE OF FLATPAK.

	// ACTION routes
	r.POST("/api/shutdown/", func(c *gin.Context) { apiShutdownServer(c) }) // shutdown the server
	r.POST("/api/stop/", func(c *gin.Context) { stopPeek(c) })              // stop the peek application

	// Serve the API
	log.Info("Verifying Peek host and port...")
	// if host equals nil, null or empty string, set to default
	if host == "" || host == "null" || host == "nil" {
		log.Warnf("Host(%s) is invalid, setting to default web UI address: %s", host, DefaultWebUiAddr)
		host = DefaultWebUiHost
	}
	// if port equals nil, null or empty string, set to default
	if port == 0 {
		log.Warnf("Port(%d) is invalid, setting to default web UI port: %d", port, DefaultWebUiPort)
		port = DefaultWebUiPort
	}

	log.Info("GIN: Peek WebServer started at address: http://" + host + ":" + strconv.Itoa(port))
	err = r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Panicf("GIN: Failed to start the Peek WebServer: %s", err)
	} // listen and serve
}
