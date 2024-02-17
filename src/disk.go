package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func storageAllDisks(c *gin.Context) {
	config, err := ConfigParser()
	if err != nil {
		log.Fatalf("STAD: Error getting config. err: %s", err)
	}
	if config.Show.ShowDisk == true {
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "This endpoint is disabled in the config.",
		})
	}

}
