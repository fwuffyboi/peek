package main

import (
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func rlKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

type RatelimitJson struct {
	Info  Info      `json:"info"`
	Debug DebugInfo `json:"debug"`
}

type Info struct {
	Error     string `json:"error"`
	UITitle   string `json:"UITitle"`
	UIMessage string `json:"UIMessage"`
}

type DebugInfo struct {
	ResetTime string `json:"resetTime"`
	Limit     int    `json:"limit"`
	ClientIP  string `json:"clientIP"`
	ServerIP  string `json:"serverIP"`
	// ServerPort string `json:"serverPort"`
}

func rlErrorHandler(c *gin.Context, info ratelimit.Info) {

	infoJson := Info{
		Error:   "ratelimit",
		UITitle: "Too many requests",
		UIMessage: "Whoops! Looks like this device is sending too many requests! Please try again in " +
			time.Until(info.ResetTime).String() +
			". This happened because the server detected this IP address making more than the configured (" +
			strconv.Itoa(int(info.Limit)) + ") requests per second.",
	}
	debugJson := DebugInfo{
		ResetTime: info.ResetTime.String(),
		Limit:     int(info.Limit),
		ClientIP:  c.ClientIP(),
		ServerIP:  strings.Split(c.Request.Host, ":")[0], // remove port from server ip
		// todo: why does the below line sometimes cause a panic if requests are made rapidly?
		// ServerPort: strings.Split(c.Request.Host, ":")[1],
	}

	RLJson := RatelimitJson{
		Info:  infoJson,
		Debug: debugJson,
	}

	c.JSON(http.StatusTooManyRequests, RLJson) // reply with above json
}
