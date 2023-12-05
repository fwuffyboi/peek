package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// Show all API endpoints
func apiEndpoints(c *gin.Context) {
	if UnsupportedOS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "This operating system is not supported. Please use a Linux or Darwin(MacOS) derivative.",
		})
	} else {
		endpoints := map[string]string{
			"/api":          "Show all API endpoints",
			"/api/full":     "Show all API stats",
			"/api/shutdown": "Shutdown the server",
		}
		// Send the JSON response
		c.JSON(http.StatusOK, gin.H{"endpoints": endpoints})
	}
}

// Show all API stats
func apiFull(c *gin.Context) {
	if UnsupportedOS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "This operating system is not supported. Please use a Linux or Darwin(MacOS) derivative.",
		})
	} else {
		uptime, err := getUptime()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
		}

		uptimeFullFriendly, uptimeFullRaw := formatUptime(uptime)
		hostname, err := os.Hostname()

		// json we shit out to the api
		c.JSON(http.StatusOK, gin.H{
			// Peek app info stuff
			"applicationName":    "Peek",
			"applicationVersion": VERSION,

			// uptime
			"uptime-seconds":           uptime.Seconds(),
			"uptime-ddhhmmss-raw":      uptimeFullRaw,
			"uptime-ddhhmmss-friendly": uptimeFullFriendly,

			// ip stuff
			"serverIP": IpAddress,
			"clientIP": c.ClientIP(),

			// country stuff
			"serverCountry": ServerCountry,
			"clientCountry": countryFromIP(c.ClientIP()),

			// hostname
			"hostname": hostname,
		})
	}
}

// Shutdown the server
func apiShutdownServer(c *gin.Context) {
	if UnsupportedOS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "This operating system is not supported. Please use a Linux or Darwin(MacOS) derivative.",
		})
	} else {
		if c.Request.Method != "POST" { // if not a posty requesty
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"err": "To interact with this API endpoint, you must use a POST request.",
			})
		} else { // if is a posty requesty westy
			if c.Query("confirm") == "true" { // if ?confirm=true in url
				// shut down server!1!!! :3
				shutdownDelay := time.Second * 85
				log.Warnf("API: Shutdown request received from client IP: %s at time: %s.",
					c.ClientIP(), time.Now().Format("2006-01-02, 15:04:05"))
				log.Warnf("API: SHUTTING DOWN IN %s SECONDS (in %d minutes)!!!", shutdownDelay, int(shutdownDelay.Minutes()))
				c.JSON(http.StatusOK, gin.H{
					"msg": c.ClientIP() + "has requested a server shutdown in " + shutdownDelay.String() + " seconds.",
				})

			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": "You must confirm the shutdown by adding ?confirm=true to the url.",
				})
			}
		}
	}
}
