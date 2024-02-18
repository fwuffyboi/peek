package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Show all API endpoints
func apiEndpoints(c *gin.Context) {
	endpoints := map[string]string{
		"GET   /":             "Show information about Peek",
		"GET   /api":          "Show all API endpoints",
		"GET   /api/full":     "Show all API stats",
		"GET   /api/logs/all": "Show all logs",
		"POST  /api/stop":     "Stop Peek",
		"POST  /api/shutdown": "Shutdown the server",
	}
	// Send the JSON response
	c.JSON(http.StatusOK, gin.H{"endpoints": endpoints})
}

type applicationStruct struct {
	ApplicationName    string `json:"applicationName"`
	ApplicationVersion string `json:"applicationVersion"`
}
type clientStruct struct {
	ClientIP      string `json:"clientIP"`
	ClientCountry string `json:"clientCountry"`
	ClientFlag    string `json:"clientFlag"`
}
type serverStruct struct {
	ServerIP      string `json:"serverIP"`
	ServerCountry string `json:"serverCountry"`
	ServerFlag    string `json:"serverFlag"`
}
type uptimeStruct struct {
	UptimeSeconds          float64 `json:"uptime-seconds"`
	UptimeDDHHMMSSRaw      string  `json:"uptime-ddhhmmss-raw"`
	UptimeDDHHMMSSFriendly string  `json:"uptime-ddhhmmss-friendly"`
}
type hostnameStruct struct {
	Hostname string `json:"hostname"`
}
type memoryStruct struct {
	MemoryTotal       uint64 `json:"memoryTotal"`
	MemoryFree        uint64 `json:"memoryFree"`
	MemoryUsed        uint64 `json:"memoryUsed"`
	MemoryUsedPercent int    `json:"memoryUsedPercent"`
}
type cpuStruct struct {
	HighestCPUTemp    string `json:"highestCPUTemp"`
	ZoneOfHighestTemp string `json:"zoneOfHighestTemp"`
	CPUUsage          string `json:"cpuUsage"`
}
type apiFullResponse struct {
	Application applicationStruct `json:"application"`
	Client      clientStruct      `json:"client"`
	Server      serverStruct      `json:"server"`
	Uptime      uptimeStruct      `json:"uptime"`
	Hostname    hostnameStruct    `json:"hostname"`
	Memory      memoryStruct      `json:"memory"`
	CPU         cpuStruct         `json:"cpu"`
}

// Show all API stats
func apiFull(c *gin.Context) {
	uptimeVar, err := getUptime()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
	}

	// create all the var shit
	var uptimeFullFriendly, uptimeFullRaw string
	var uptimeSeconds float64
	var hostnameVar string
	var clientCountry string
	var serverIP string
	var clientFlag string
	var serverFlag string
	var memoryTotal, memoryFree, memoryUsed uint64
	var memoryUsedPercent int
	var HCPUTemp string
	var HCPUZone string
	var CPUUse string

	config, err := ConfigParser()

	if config.Show.ShowUptime == false {
		uptimeFullFriendly = "This value is disabled."
		uptimeFullRaw = "This value is disabled."
		uptimeSeconds = 0
	} else {
		uptimeFullFriendly, uptimeFullRaw = formatUptime(uptimeVar)
		uptimeSeconds = uptimeVar.Seconds()
	}

	if config.Show.ShowHostname == false {
		hostnameVar = "This value is disabled."
	} else {
		hostnameVar, err = os.Hostname()
	}

	if config.Show.ShowClientCountry == false {
		clientCountry = "This value is disabled."
		clientFlag = "This value is disabled."
	} else {
		cip := c.ClientIP()
		if cip == "127.0.0.1" || cip == "0.0.0.0" {
			clientFlag = "LOCALHOST"
			clientCountry = "LOCALHOST"
		} else {
			clientCountry = countryFromIP(cip)
			clientFlag = "https://flagpedia.net/data/flags/emoji/twitter/256x256/" + clientCountry + ".png"
		}

	}

	if config.Show.ShowServerCountry == false {
		ServerCountry = "This value is disabled."
		serverFlag = "This value is disabled."
	} else {
		ServerCountry = countryFromIP(IpAddress)
		serverFlag = "https://flagpedia.net/data/flags/emoji/twitter/256x256/" + strings.ToLower(ServerCountry) + ".png"
	}

	if config.Show.ShowIP == false {
		serverIP = "This value is disabled."
	} else {
		serverIP = IpAddress
	}

	if config.Show.ShowRAM == false {
		memoryTotal, memoryFree, memoryUsed, memoryUsedPercent = 0, 0, 0, 0
	} else {
		memoryTotal, memoryFree, memoryUsed, memoryUsedPercent, err = getMemoryUsage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
		}
	}

	if config.Show.ShowCPUTemp == false {
		HCPUTemp = "This value is disabled."
		HCPUZone = "This value is disabled."
	} else {
		HCPUTemp, HCPUZone, err = GetHighestCPUTemp()
		if HCPUTemp == "ERROR" || HCPUZone == "UNKNOWN" {
			HCPUTemp = "ERROR"
			HCPUZone = "ERROR"
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
		}
	}
	if config.Show.ShowCPUUsage == false {
		CPUUse = "This value is disabled."
	} else {
		CPUUse, err = GetCPUUsage()
		if CPUUse == "ERROR" {
			CPUUse = "ERROR"
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
		}
	}

	// Send the JSON response
	c.JSON(http.StatusOK, apiFullResponse{
		Application: applicationStruct{
			ApplicationName:    "Peek",
			ApplicationVersion: VERSION,
		},
		Client: clientStruct{
			ClientIP:      c.ClientIP(),
			ClientCountry: clientCountry,
			ClientFlag:    clientFlag,
		},
		Server: serverStruct{
			ServerIP:      serverIP,
			ServerCountry: ServerCountry,
			ServerFlag:    serverFlag,
		},
		Uptime: uptimeStruct{
			UptimeSeconds:          uptimeSeconds,
			UptimeDDHHMMSSRaw:      uptimeFullRaw,
			UptimeDDHHMMSSFriendly: uptimeFullFriendly,
		},
		Hostname: hostnameStruct{
			Hostname: hostnameVar,
		},
		Memory: memoryStruct{
			MemoryTotal:       memoryTotal,
			MemoryFree:        memoryFree,
			MemoryUsed:        memoryUsed,
			MemoryUsedPercent: memoryUsedPercent,
		},
		CPU: cpuStruct{
			HighestCPUTemp:    HCPUTemp,
			ZoneOfHighestTemp: HCPUZone,
			CPUUsage:          CPUUse,
		},
	})

}

// Shutdown the server
func apiShutdownServer(c *gin.Context) {
	config, err := ConfigParser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
	}
	if config.Actions.SystemShutdown == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "This endpoint is disabled in the config.",
		})
	} else {
		if c.Request.Method != "POST" { // if not a posty requesty
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"err": "To interact with this API endpoint, you must use a POST request.",
			})
		} else { // if is a posty requesty westy
			if c.Query("confirm") == "true" { // if ?confirm=true in url
				// shut down server!1!!! :3
				shutdownDelay := config.Api.ShutdownDelay
				log.Infof("API: Shutdown request received from client IP: %s at time: %s. Waiting with a delay of %d minutes until shutdown.",
					c.ClientIP(), time.Now().Format("2006-01-02, 15:04:05"), shutdownDelay)
				time.Sleep(time.Duration(shutdownDelay))
				minArg := "+" + strconv.Itoa(shutdownDelay)
				cmd := exec.Command("shutdown", "-P", minArg)

				outputBytes, err := cmd.CombinedOutput()
				if err != nil {
					log.Errorf("Error shutting down: %s", err)
				} else {
					fmt.Println(string(outputBytes))
					c.JSON(http.StatusOK, gin.H{
						"msg":        c.ClientIP() + " has requested a server shutdown in " + strconv.Itoa(shutdownDelay) + " minutes.",
						"cmd_output": string(outputBytes),
					})
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
func stopPeek(c *gin.Context) { // TODO: add auth
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

// Return the logs
func apiLogs(c *gin.Context) {
	config, _ := ConfigParser()
	if config.Show.ShowLogsAPI == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "This endpoint is disabled in the config.",
		})
		return
	}
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
