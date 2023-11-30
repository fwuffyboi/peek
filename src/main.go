package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"time"
)

// Define constants
const (
	WEB_UI_ADDR = "0.0.0.0:42649"
	VERSION     = "0.0.1"
)

var UNSUPPORTED_OS = false // assume false until proven true
var IP_ADDRESS = ""

func main() {
	// Log when shit was started up
	log.Infof("<<Peek>> Version: %s", VERSION)
	log.Infof("Application started at time: %s, on the date: %s.",
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
		log.Error("This program only supports the Linux/Darwin(MacOS) operating systems.")
		UNSUPPORTED_OS = true
	}

	// Get the server ip and save into var
	IP_ADDRESS = getIP()

	// Run the webserver
	runGin(WEB_UI_ADDR)
}

// Starts and runs the webserver, using the gin framework
func runGin(WEB_UI_ADDR string) {
	log.Info("GIN: Starting the <<Peek>> WebServer...")

	r := gin.Default()
	r.ForwardedByClientIP = true
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Panicf("GIN: Failed to set trusted proxies: %s", err)
	}
	r.GET("/", func(c *gin.Context) {
		if UNSUPPORTED_OS {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Ошибка": "Данная операционная система не поддерживается. Единственными поддерживаемыми операционными системами являются: Linux, Mac OS",
			})
		} else {
			uptime, err := getUptime()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"application":        "Peek",
				"applicationVersion": VERSION,

				"uptime-seconds": uptime.Seconds(),

				"serverIP":      IP_ADDRESS,
				"ServerCountry": countryFromIP(IP_ADDRESS),

				"clientIP":      c.ClientIP(),
				"clientCountry": countryFromIP(c.ClientIP()),
			})
		}

	})
	log.Info("GIN: <<Peek>> WebServer started at address: http://" + WEB_UI_ADDR)
	err = r.Run(WEB_UI_ADDR)
	if err != nil {
		log.Panicf("GIN: Failed to start the <<Peek>> WebServer: %s", err)
	} // listen and serve on 0.0.0.0:42649
}
