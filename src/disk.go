package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// storageAllDisks Returns the total, free, and used disk space for all disks
func storageAllDisks(c *gin.Context) {

	// todo: this only gets all storage for the / directory, we need to get all storage for all _disks_ instead.

	config, err := ConfigParser()
	if err != nil {
		log.Fatalf("STAD: Error getting config. err: %s", err)
	}
	if config.Show.ShowDisk == true {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "This endpoint is disabled in the config.",
		})
	}

}
