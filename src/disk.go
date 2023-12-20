package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/disk"
	log "github.com/sirupsen/logrus"
)

func storageAllDisks(c *gin.Context) {
	// Get the disk usage for the root directory ("/" on Unix-like systems)
	usage, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	UP3 := fmt.Sprintf("%.3f", usage.UsedPercent)
	c.JSON(200, gin.H{
		"disk_total":        usage.Total,
		"disk_free":         usage.Free,
		"disk_used":         usage.Used,
		"disk_used_percent": UP3,
	})
}
