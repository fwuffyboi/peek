package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetHighestCPUTemp Gets the highest temperature of a cpu thing
func GetHighestCPUTemp() (string, string, error) {
	var thermalZonePath = "/sys/class/thermal"
	var highestTemp = 0
	var highestZone string
	var cTemp int
	var cZone string

	files, err := os.ReadDir(thermalZonePath)
	if err != nil {
		log.Error("[GHCT]   Failed to read thermalZonePath directory. Err: ", err)
		return "ERROR", "UNKNOWN", err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "thermal_zone") {
			zonePath := filepath.Join(thermalZonePath, file.Name())        // join file name together
			tempBytes, err := os.ReadFile(filepath.Join(zonePath, "temp")) // read temp
			if err != nil {                                                // if error, stop
				return "ERROR", "UNKNOWN", err
			}

			temp, err := strconv.Atoi(strings.TrimSpace(string(tempBytes))) // make temp an int
			if err != nil {                                                 // if error, stop
				log.Error("[GHCT]   Failed to convert temp to int. Err: ", err) // ;-;
				return "ERROR", "UNKNOWN", err
			}

			cTemp = temp / 1000 // Convert milli-degrees to degrees Celsius
			cZone = file.Name()
			if highestTemp < cTemp {
				highestTemp = cTemp
				highestZone = cZone
			}
			log.Infof("[GHCT]   cTemp: %d, cZone: %s", cTemp, cZone)
		}
	}

	return strconv.Itoa(highestTemp), highestZone, nil
}

func GetCPUUsage() (string, error) {
	// Get CPU usage percentages for all CPU cores
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		log.Error("[GCPU]   Error:", err)
		return "ERROR", err
	}

	// Calculate total CPU usage
	totalUsage := 0.0
	for _, percent := range percentages {
		totalUsage += percent
	}

	// Calculate average CPU usage
	averageUsage := totalUsage / float64(len(percentages))

	// Print average CPU usage
	log.Infof("[GCPU]   Average CPU Usage: %.2f%%", averageUsage)
	return fmt.Sprintf("%.2f", averageUsage), nil
}
