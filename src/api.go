package main

import (
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// constants
const disabledValueText = "This value is disabled."

// Show all API endpoints
func apiEndpoints(c *gin.Context) {
	endpoints := map[string]string{
		"GET   /":                         "Show information about Peek",
		"GET   /api/":                     "Show all API endpoints",
		"GET   /api/stats/all/":           "Show all API stats",
		"GET   /api/logs/all/":            "Show all logs",
		"GET   /api/heartbeat/":           "Show if Peek is online/responsive",
		"POST  /api/stop/":                "Stop Peek",
		"POST  /api/shutdown/":            "Shutdown the server",
		"GET   /api/alerts/":              "Show all alerts",
		"POST  /api/auth/create/session/": "Create an auth token",
		"POST  /api/auth/verify/session/": "Verify an auth token",
	}
	// Send the JSON response
	c.JSON(http.StatusOK, gin.H{"endpoints": endpoints})
}

type applicationStruct struct {
	Name    string `json:"name"`
	Version string `json:"version"`
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
	ServerTZ      string `json:"serverTimezone"`
	ServerTime    string `json:"serverTime"`
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
	MemoryTotal       int `json:"memoryTotal"`
	MemoryFree        int `json:"memoryFree"`
	MemoryUsed        int `json:"memoryUsed"`
	MemoryUsedPercent int `json:"memoryUsedPercent"`

	SwapUsed        int `json:"swapUsed"`
	SwapTotal       int `json:"swapTotal"`
	SwapPercentUsed int `json:"swapPercentUsed"`
}
type cpuStruct struct {
	HighestCPUTemp    string `json:"highestCPUTemp"`
	ZoneOfHighestTemp string `json:"zoneOfHighestTemp"`
	CPUUsage          string `json:"usage"`

	CPUVendor    string `json:"vendor"`
	CPUModel     string `json:"model"`
	CPUModelName string `json:"modelName"`
	CPUCores     int    `json:"cores"`
	CPUMhz       int    `json:"mhz"`
	CPUCacheSize int    `json:"cacheSize"`
}
type diskStruct struct {
	Disks []DiskUsage `json:"disks"`
}
type apiFullResponse struct {
	Application applicationStruct    `json:"application"`
	Client      clientStruct         `json:"client"`
	Server      serverStruct         `json:"server"`
	Uptime      uptimeStruct         `json:"uptime"`
	Hostname    hostnameStruct       `json:"hostname"`
	Memory      memoryStruct         `json:"memory"`
	CPU         cpuStruct            `json:"cpu"`
	Disk        diskStruct           `json:"disk"`
	Alerts      map[string]time.Time `json:"alerts"`
}

// @Summary Shows every statistic on the server that is enabled in the configuration
// @Description  Shows every statistic on the server that is enabled in the configuration. This includes the server's IP, the server's country, the server's timezone, the server's time, the server's hostname, the server's uptime, the server's memory usage, the server's swap usage, the server's CPU usage, the server's CPU temperature, the server's CPU model, the server's CPU vendor, the server's CPU model name, the server's CPU cores, the server's CPU MHz, the server's CPU cache size, the client's IP, the client's country, the client's flag, and the server's flag.
// @Accept  json
// @Produce  json
// @Success 200 "All stats are shown in JSON format"
// @Failure 401 "Unauthorized/auth error, make a session first and use that token to access this endpoint"
// @Failure 429
// @param Authorization header string false "The auth token to use to access this endpoint."
// @Tags apiStatsGroup
// @Router /stats/all [get]
// allStatsAPI Shows all API stats
func allStatsAPI(c *gin.Context) {

	// check if user requires auth
	config, err := ConfigParser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return // terminate completely
	}

	if config.Auth.AuthRequired {
		if !isAuthed(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Unauthorized",
				"msg": "You must be authenticated to access this endpoint.",
			})
			return
		}
	}

	uptimeVar, err := getUptime()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
	}

	// server uptime
	var uptimeFullFriendly, uptimeFullRaw string
	var uptimeSeconds float64

	// hostname
	var hostnameVar string

	// client networking
	var clientCountry string
	var clientFlag string

	// server networking
	var serverIP string
	var serverFlag string

	// memory/ram
	var memoryTotal, memoryFree, memoryUsed int
	var memoryUsedPercent int

	// swap
	var swapUse, swapTotal, swapPercentUsed int

	// cpu
	var HCPUTemp string
	var HCPUZone string
	var CPUUse string

	// CPU info
	var CPUVendor, CPUModel, CPUModelName string
	var CPUCores, CPUMhz, CPUCacheSize int

	// disks
	var DiskUsageList []DiskUsage

	// time
	var serverTZ string
	var serverTime string

	if !config.Show.ShowUptime {
		uptimeFullFriendly = disabledValueText
		uptimeFullRaw = disabledValueText
		uptimeSeconds = 0
	} else {
		uptimeFullFriendly, uptimeFullRaw = formatUptime(uptimeVar)
		uptimeSeconds = uptimeVar.Seconds()
	}

	if !config.Show.ShowHostname {
		hostnameVar = disabledValueText
	} else {
		hostnameVar, err = os.Hostname()
		if err != nil {
			hostnameVar = "Error."
		}
	}

	if !config.Show.ShowClientCountry {
		clientCountry = disabledValueText
		clientFlag = disabledValueText
	} else {
		cip := c.ClientIP()
		if cip == "127.0.0.1" || cip == "0.0.0.0" || cip == "::1" {
			clientFlag = "unknown"
			clientCountry = "unknown"
		} else {
			clientCountry = countryFromIP(cip)
			clientFlag = "https://flagpedia.net/data/flags/emoji/twitter/256x256/" + strings.ToLower(clientCountry) + ".png"
		}

	}

	if !config.Show.ShowServerCountry {
		ServerCountry = disabledValueText
		serverFlag = disabledValueText
	} else {
		ServerCountry = countryFromIP(ServerIPAddress)
		serverFlag = "https://flagpedia.net/data/flags/emoji/twitter/256x256/" + strings.ToLower(ServerCountry) + ".png"
	}

	if !config.Show.ShowIP {
		serverIP = disabledValueText
	} else {
		serverIP = ServerIPAddress
	}
	if !config.Show.ShowTimezone {
		serverTZ = disabledValueText
		serverTime = disabledValueText
	} else {
		serverTZ = serverTimezone()
		serverTime = time.Now().Format("2006-01-02, 15:04:05")
	}

	if !config.Show.ShowRAM {
		memoryTotal, memoryFree, memoryUsed, memoryUsedPercent = 0, 0, 0, 0
		swapPercentUsed, swapUse, swapTotal = 0, 0, 0
	} else {
		memoryTotal, memoryFree, memoryUsed, memoryUsedPercent, err = getMemoryUsage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
		}

		// get swap too
		swapUse, swapTotal, swapPercentUsed, err = getSwapUsage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})

		}
	}

	if !config.Show.ShowCPUTemp {
		HCPUTemp = disabledValueText
		HCPUZone = disabledValueText
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
	if !config.Show.ShowCPUUsage {
		CPUUse = disabledValueText
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

	if !config.Show.ShowCPU {
		CPUVendor, CPUModel, CPUModelName, CPUCores, CPUMhz, CPUCacheSize = disabledValueText, disabledValueText, disabledValueText, 0, 0, 0
	} else {
		CPUVendor, CPUModel, CPUModelName, CPUCores, CPUMhz, CPUCacheSize, err = getCPUInfo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
		}
		log.Infof("CPU: Vendor: %s, Model: %s, Model Name: %s, Cores: %d, MHz: %d, Cache Size: %d", CPUVendor, CPUModel, CPUModelName, CPUCores, CPUMhz, CPUCacheSize)
	}

	if !config.Show.ShowDisk {
		DiskUsageList = nil
	} else {
		DiskUsageList, err = storageAllDisks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
			return
		}

	}

	// Send the JSON response
	c.JSON(http.StatusOK, apiFullResponse{
		Application: applicationStruct{
			Name:    "Peek",
			Version: VERSION,
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
			ServerTZ:      serverTZ,
			ServerTime:    serverTime,
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

			SwapUsed:        swapUse,
			SwapTotal:       swapTotal,
			SwapPercentUsed: swapPercentUsed,
		},
		CPU: cpuStruct{
			HighestCPUTemp:    HCPUTemp,
			ZoneOfHighestTemp: HCPUZone,
			CPUUsage:          CPUUse,

			CPUVendor:    CPUVendor,
			CPUModel:     CPUModel,
			CPUModelName: CPUModelName,
			CPUCores:     CPUCores,
			CPUMhz:       CPUMhz,
			CPUCacheSize: CPUCacheSize,
		},
		Disk: diskStruct{
			Disks: DiskUsageList,
		},
		Alerts: getAlerts(),
	})

}

// @Summary Shuts down the server
// @Description Shuts down the server
// @Accept  json
// @Produce  json
// @Success 200 "All stats are shown in JSON format"
// @Failure 400 "This endpoint is disabled in the config."
// @Failure 401 "Unauthorized/auth error, make a session first and use that token to access this endpoint"
// @Failure 429
// @Failure 500
// @param Authorization header string false "The auth token to use to access this endpoint"
// @Tags apiPeekGroup
// @Router /peek/shutdown [put]
// apiShutdownServer Shuts down the server
func apiShutdownServer(c *gin.Context) {
	config, err := ConfigParser()
	if err != nil {
		log.Errorf("Failed to get config: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	if !config.Actions.SystemShutdown {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "This endpoint is disabled in the config.",
		})
		return
	}

	if config.Auth.AuthRequired {
		if !isAuthed(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Unauthorized",
				"msg": "You must be authenticated to access this endpoint.",
			})
			return
		}
	}

	shutdownDelay := config.Api.ShutdownDelay
	log.Infof("API: Shutdown request received from client IP: %s at time: %s. Waiting with a delay of %d minutes until shutdown.",
		c.ClientIP(), time.Now().Format("2006-01-02, 15:04:05"), shutdownDelay)
	time.Sleep(time.Duration(shutdownDelay))
	minArg := "+" + strconv.Itoa(shutdownDelay)
	cmd := exec.Command("shutdown", "-P", minArg)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Error shutting down: %s", err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":        c.ClientIP() + " has requested a server shutdown in " + strconv.Itoa(shutdownDelay) + " minutes.",
			"cmd_output": string(outputBytes),
		})
		return
	}
}

// @Summary Stop the Peek application
// @Description Stop the Peek application
// @Accept  json
// @Produce  json
// @Success 200 "All stats are shown in JSON format"
// @Failure 400 "This endpoint is disabled in the config."
// @Failure 429
// @Failure 500 "Internal Server Error"
// @Tags apiPeekGroup
// @param Authorization header string false "The auth token to use to access this endpoint"
// @Router /peek/stop [put]
// stopPeek Stops the Peek application
func stopPeek(c *gin.Context) {
	config, err := ConfigParser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to get config: " + err.Error(),
		})
	}

	if !config.Actions.ShutdownPeek {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "This endpoint is disabled in the config.",
		})
		return
	}

	if config.Auth.AuthRequired {
		if !isAuthed(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Unauthorized",
				"msg": "You must be authenticated to access this endpoint.",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": c.ClientIP() + " has requested that Peek stops. Stopping this application NOW!",
	}) // TODO: make this actually respond to client

	log.Warnf("SHUTDOWN: %s has made a Peek shutdown request.", c.ClientIP())
	log.Warn("Peek is shutting down...")
	log.Fatalf("Peek has been shut down due to a client's request. Client's info: IP: %s, Country: %s", c.ClientIP(), countryFromIP(c.ClientIP()))
	return
}

// @Summary Returns the logfile in plain text
// @Description Returns the logfile in plain text
// @Accept  json
// @Produce  json
// @Success 200 "The logfile is returned in plain text"
// @Failure 400 "This endpoint is disabled in the config."
// @Failure 429
// @Failure 500 "Internal Server Error"
// @param download query bool false "If true, the logfile will be downloaded"
// @param Authorization header string false "The auth token to use to access this endpoint"
// @Tags apiLogsGroup
// @Router /logs/all [get]
// Return the logs
func apiLogs(c *gin.Context) {
	config, err := ConfigParser()
	if err != nil {
		log.Errorf("Failed to get config: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	if !config.Show.ShowLogsAPI {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "This endpoint is disabled in the config.",
		})
		return
	}

	if config.Auth.AuthRequired {
		if !isAuthed(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Unauthorized",
				"msg": "You must be authenticated to access this endpoint.",
			})
			return
		}
	}

	if c.Query("download") == "true" {
		usrHome, _ := os.UserHomeDir()
		peekLogPath := path.Join(usrHome, ".config/peek", "peek.log")
		c.FileAttachment(peekLogPath, "peek.log")
		return
	} else {
		usrHome, _ := os.UserHomeDir()
		peekLogPath := path.Join(usrHome, ".config/peek", "peek.log")
		fileContents, err := os.ReadFile(peekLogPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err,
			})
		}
		c.String(200, string(fileContents))
		return
	}
}

// @Summary Shows a page with a very simple heartbeat message
// @Description Shows a page with a very simple heartbeat message saying "online" to show that the server is online and responding to requests.
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 429
// @Tags apiInfoGroup
// @Router /heartbeat [get]
func apiHeartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "online",
	})
}
