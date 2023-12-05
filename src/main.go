package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

// Define constants
const (
	WebUiAddr = "0.0.0.0:42649" // Address of the webserver, HAS to be in the format of: IP:PORT
	VERSION   = "0.0.1"         // Version of Peek
)

var UnsupportedOS = false // assume false until proven true
var IpAddress = ""        // IP address of the server
var ServerCountry = ""    // Country of the server, based on IP

func main() {
	// Log when shit was started up
	log.Infof("<<Peek>> Version: %s", VERSION)
	log.Infof("Application started at time: %s, on the date: %s(YYYY-MM-DD).",
		time.Now().Format("15:04:05"),   // Format for hh:mm:ss
		time.Now().Format("2006-01-02"), // Format for yyyy-mm-dd
	)

	// Copyright notices
	log.Info("<<Peek>> is licensed under the MIT License. See LICENSE for more information.")
	log.Infof("(C) %s <<Peek>> Contributors. All rights reserved.", time.Now().Format("2006"))
	log.Info("<<Peek>> is a program written by: @fwuffyboi (https://github.com/fwuffyboi)")
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

	// Run the webserver
	runGin(WebUiAddr)
}

// Starts and runs the webserver, using the gin framework
func runGin(WebUiAddr string) {
	log.Info("GIN: Starting the <<Peek>> WebServer...")

	gin.SetMode(gin.ReleaseMode) // set to production mode
	r := gin.Default()
	r.ForwardedByClientIP = true

	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Panicf("GIN: Failed to set trusted proxies: %s", err)
	}

	r.GET("/api/", func(c *gin.Context) { apiEndpoints(c) })                                    // return all api endpoints
	r.GET("/api/full/", func(c *gin.Context) { apiFull(c) })                                    // return all api/json info
	r.POST("/api/shutdown/", shutdownMiddleware, func(c *gin.Context) { apiShutdownServer(c) }) // shutdown the server
	r.POST("/api/stop/", func(c *gin.Context) { shutdownPeek(c) })                              // stop the peek application

	log.Info("GIN: <<Peek>> WebServer started at address: http://" + WebUiAddr)
	err = r.Run(WebUiAddr)
	if err != nil {
		log.Panicf("GIN: Failed to start the <<Peek>> WebServer: %s", err)
	} // listen and serve on 0.0.0.0:42649
}
