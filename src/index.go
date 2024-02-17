package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func indexPage(c *gin.Context) {

	var content = "Hello, this is device is running Peek, the easy-to-use system monitoring tool. To view the API endpoints, go to /api. If you would like to use the WebUI to view your server's statistics, go to https://example.com."
	var GHLink = "https://github.com/fwuffyboi/peek"
	var copyright = "(C) " + time.Now().Format("2006") + " fwuffyboi(Эшли Карамель/Ashley Caramel) and contributors."

	// Send the text response
	c.JSON(http.StatusOK, gin.H{
		"UIContent": content,
		"GitHub":    GHLink,
		"Copyright": copyright,
	})
}
