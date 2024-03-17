package main

import (
	"fmt"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/fwuffyboi/peek/src/docs"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginlogrus "github.com/toorop/gin-logrus"
	"net/http"
	"strconv"
	"time"
)

// Comments for swagger to build docs from
// @title Swagger for the Peek API
// @version v0.9.0

// @contact.name API Support
// @contact.url https://github.com/fwuffyboi/peek/issues

// @license.name MIT
// @license.url https://github.com/fwuffyboi/peek/blob/main/LICENSE

// runGin Starts and runs the webserver, using the gin framework
func runGin(host string, port int, ginRatelimit int) {

	// log
	log.Info("GIN: Starting the Peek WebServer...")

	// Set up rate limiter
	log.Info("GIRL: Setting up rate limiter...")
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

	// This middleware adds cors headers and informational headers
	r.Use(func(c *gin.Context) {

		// cors headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// if a req is an options req, just return 200, as this is only used for cors in this project
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// informational headers
		c.Writer.Header().Set("Server", "Gin")

		// continue to next middleware
		c.Next()
	})
	r.Use(ginlogrus.Logger(log.StandardLogger()), rl)

	// Set up trusted proxies, so we can get the real ip from the client, if behind a proxy
	r.ForwardedByClientIP = true // set to true, so we can get the real ip from the client

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

	// Define the routes
	docs.SwaggerInfo.BasePath = "/api/v1" // set the base path for the swagger docs
	v1API := r.Group("/api/v1")
	{
		apiInfoGroup := v1API.Group("/") // information only regarding the project or api itself
		{
			apiInfoGroup.GET("/", rl, func(c *gin.Context) { indexPage(c) })          // return the web ui
			apiInfoGroup.GET("/heartbeat/", func(c *gin.Context) { apiHeartbeat(c) }) // return an "online" message

			// apiInfoGroup.GET("/api/", rl, func(c *gin.Context) { apiEndpoints(c) }) // return all api endpoints todo: should i remove this?
		}

		apiStatsGroup := v1API.Group("/stats") // stats/info about the server itself
		{
			apiStatsGroup.GET("/all/", rl, func(c *gin.Context) { allStatsAPI(c) }) // Return all api/json info
		}

		apiLogsGroup := v1API.Group("/logs") // logs
		{
			apiLogsGroup.GET("/all/", rl, func(c *gin.Context) { apiLogs(c) }) // return everything in the logfile
			// todo: be able to get logs from a particular event, such as: auth, ip2country, etc
		}

		apiAuthGroup := v1API.Group("/auth") // authentication
		{
			apiAuthGroup.POST("/create/session/", rl, func(c *gin.Context) { createSession(c) }) // create an auth token
			// apiAuthGroup.POST("/verify/session/", rl, func(c *gin.Context) { verifySession(c) }) // verify an auth token
		}

		apiPeekGroup := v1API.Group("/peek") // peek
		{
			apiPeekGroup.POST("/stop/", rl, func(c *gin.Context) { stopPeek(c) })              // stop the peek application
			apiPeekGroup.POST("/shutdown/", rl, func(c *gin.Context) { apiShutdownServer(c) }) // shutdown the server
			apiPeekGroup.GET("/alerts/", rl, func(c *gin.Context) { apiReturnAlerts(c) })      // return all alerts
		}
	}

	// swagger docs
	if config.Api.SwaggerEnabled == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	} else {
		log.Warn("Swagger is disabled, so the docs will not be available")
		r.GET("/swagger/*any", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "Swagger is disabled, to turn on, set swagger-enabled to true in the config file's 'api' section and restart the application."})
		})
	}

	// permanent redirect to /api/v1
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/api/v1") })

	// Serve the API
	log.Info("Verifying Peek host and port...")
	// if host equals nil, null or empty string, set to default
	if host == "" || host == "null" || host == "nil" {
		log.Warnf("Host(%s) is invalid, setting to default web UI address: %s", host, DefaultWebuiAddress)
		host = DefaultWebUiHost
	}
	// if port equals 0 or less than 1025, set to default
	if port == 0 || port <= 1024 {
		log.Warnf("Port(%d) is invalid, setting to default web UI port: %d", port, DefaultWebUiPort)
		port = DefaultWebUiPort
	}

	log.Info("GIN: Peek WebServer started at address: http://" + host + ":" + strconv.Itoa(port))
	err = r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("GIN: Failed to start the Peek WebServer: %s", err)
	}
}
