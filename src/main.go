package main

import (
	"fmt"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Define constants
const (
	// DefaultWebuiAddress Default address for the web UI, in case it is not provided/invalid in the config file
	DefaultWebuiAddress = "0.0.0.0:42649" // Address of the webserver, HAS to be in the format of: IP:PORT

	// VERSION Version of Peek
	VERSION = "v0.9.0-alpha" // Version of Peek

	// DefaultWebUiHost Default host for the web UI, in case it is not provided/invalid in the config file
	DefaultWebUiHost = "0.0.0.0"
	// DefaultWebUiPort Default port for the web UI, in case it is not provided/invalid in the config file
	DefaultWebUiPort = 42649
)

var ServerIPAddress = "" // IP address of the server
var ServerCountry = ""   // Country of the server, based on IP

var alertsList = make(map[string]time.Time) // List of alerts

func main() {
	// Setup logging and obtain the log file handle and multi-writer
	logFile, _, err := setupLogging()
	if err != nil {
		panic(err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			panic(err)
		}
	}(logFile)

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

	// Check os, if windows/darwin, panic, else continue
	log.Info("Checking what operating system is in use...")
	switch runtime.GOOS {
	case "linux":
		log.Info("Linux derivative detected.")
	default:
		log.Error("Unsupported operating system detected.")
		log.Fatal("This program only supports Linux distributions.")
	}

	// start the update checker thread
	go CheckForPeekUpdate()

	// Get the server ip and save into var
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

// Starts and runs the webserver, using the gin framework
func runGin(host string, port int, ginRatelimit int) {
	log.Info("GIN: Starting the Peek WebServer...")

	log.Info("GIRL: Setting up rate limiter...")
	// Create a limiter
	if ginRatelimit == 0 {
		log.Warnf("GIRL: Rate limit is set to unlimited, setting ratelimit to 1000")
		ginRatelimit = 1000
	}
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: uint(ginRatelimit),
	})
	rl := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: rlErrorHandler,
		KeyFunc:      rlKeyFunc,
	})

	gin.SetMode(gin.ReleaseMode) // set to production mode
	r := gin.Default()           // create a gin router

	// cors middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})
	r.Use(ginlogrus.Logger(log.StandardLogger()), rl)

	r.ForwardedByClientIP = true

	// Set up trusted proxies
	config, err := ConfigParser()
	if err != nil {
		log.Fatalf("Failed to get config: %s", err)
	}

	log.Infof("Trusted proxies: %v", config.Api.TrustedProxies)
	err = r.SetTrustedProxies(config.Api.TrustedProxies)
	if err != nil {
		log.Fatalf("GIN: Failed to set trusted proxies: %s", err)
	}

	// Define API paths
	// Routes that are purely informational
	r.GET("/", rl, func(c *gin.Context) { indexPage(c) })                  // return the web ui
	r.GET("/api/", rl, func(c *gin.Context) { apiEndpoints(c) })           // return all api endpoints
	r.GET("/api/heartbeat/", func(c *gin.Context) { apiHeartbeat(c) })     // return all stats
	r.GET("/api/alerts/", rl, func(c *gin.Context) { apiReturnAlerts(c) }) // return all stats

	// Routes that cannot take user input
	r.GET("/api/stats/all/", rl, func(c *gin.Context) { apiFull(c) }) // Return all api/json info
	r.GET("/api/logs/all/", rl, func(c *gin.Context) { apiLogs(c) })  // return everything in the logfile

	// Routes that requite user input
	r.POST("/api/shutdown/", rl, func(c *gin.Context) { apiShutdownServer(c) }) // shutdown the server
	r.POST("/api/stop/", rl, func(c *gin.Context) { stopPeek(c) })              // stop the peek application

	// Serve the API
	log.Info("Verifying Peek host and port...")
	// if host equals nil, null or empty string, set to default
	if host == "" || host == "null" || host == "nil" {
		log.Warnf("Host(%s) is invalid, setting to default web UI address: %s", host, DefaultWebuiAddress)
		host = DefaultWebUiHost
	}
	// if port equals 0 or under 1025, set to default
	if port == 0 || port < 1025 {
		log.Warnf("Port(%d) is invalid, setting to default web UI port: %d", port, DefaultWebUiPort)
		port = DefaultWebUiPort
	}

	log.Info("GIN: Peek WebServer started at address: http://" + host + ":" + strconv.Itoa(port))
	err = r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("GIN: Failed to start the Peek WebServer: %s", err)
	}
}
