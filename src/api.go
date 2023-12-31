package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
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
			"GET   /api":          "Show all API endpoints",
			"GET   /api/full":     "Show all API stats",
			"POST  /api/shutdown": "Shutdown the server",
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
				time.Sleep(shutdownDelay)
				cmd := exec.Command("shutdown", "-h", "now")

				err := cmd.Run()
				if err != nil {
					log.Errorf("Error shutting down: %s", err)
				}

			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": "You must confirm the shutdown by adding ?confirm=true to the url.",
				})
			}
		}
	}
}

// Shutdown peek
func stopPeek(c *gin.Context) {
	if UnsupportedOS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "This operating system is not supported. Please use a Linux or Darwin(MacOS) derivative.",
		})
	} else {
		if c.Request.Method != "POST" { // if not a post request
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"err": "To interact with this API endpoint, you must use a POST request.",
			})
		} else { // if is a post request
			if c.Query("confirm") == "true" { // if ?confirm=true in url
				defer func() {
					log.Warnf("SHUTDOWN: %s has made a Peek shutdown request.", c.ClientIP())
					log.Warn("Peek is shutting down...")
					os.Exit(0)
				}()

				c.JSON(http.StatusOK, gin.H{
					"msg": c.ClientIP() + " has requested that Peek shuts down. Shutting down NOW!",
				}) // TODO: make this actually respond to client

			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": "You must confirm the application shutdown by adding ?confirm=true to the url.",
				})
			}

		}
	}
}

// Return the logs
func apiLogs(c *gin.Context) { // TODO: add auth
	if UnsupportedOS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "This operating system is not supported. Please use a Linux or Darwin(MacOS) derivative.",
		})
	} else {
		if c.Query("download") == "true" {
			c.FileAttachment("peek.log", "peek.log")
			return
		} else {
			fileContents, err := os.ReadFile("peek.log")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"err": err,
				})
			}
			c.String(200, string(fileContents))
		}

	}
}
