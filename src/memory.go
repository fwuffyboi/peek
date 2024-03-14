package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"strconv"
	"strings"
)

func getMemoryUsage() (total int, free int, used int, usedPercent int, err error) {

	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, 0, err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			_, err = fmt.Sscanf(fields[1], "%d", &total)
			if err != nil {
				log.Errorf("Error getting total memory: %s", err)
				return 0, 0, 0, 0, err
			}
		case "MemFree:":
			_, err = fmt.Sscanf(fields[1], "%d", &free)
			if err != nil {
				log.Errorf("Error getting free memory: %s", err)
				return 0, 0, 0, 0, err
			}
		case "MemAvailable:":
			_, err = fmt.Sscanf(fields[1], "%d", &free)
			if err != nil {
				log.Errorf("Error getting available memory: %s", err)
				return 0, 0, 0, 0, err
			}
		}
	}

	used = total - free
	percentUsed := int(math.Round(float64(used) / float64(total) * 100))

	return total, free, used, percentUsed, nil
}

func getSwapUsage() (swapUse int, swapTotal int, swapPercentUsed int, err error) {
	// Swap usage

	content, err := os.ReadFile("/proc/swaps")
	if err != nil {
		log.Errorf("Error reading /proc/swaps: %s", err)
		return 0, 0, 0, err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "Filename" {
			continue
		}
		filename := fields[0]
		total, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			log.Errorf("Error parsing used memory from /proc/swaps: %s", err)
			return 0, 0, 0, err
		}
		used, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			log.Errorf("Error parsing total memory from /proc/swaps: %s", err)
			return 0, 0, 0, err
		}

		percentUsed := int(math.Round(float64(used) / float64(total) * 100))
		log.Infof("Swap usage: %s: %dMB/%dMB (%d%%)", filename, used/1024, total/1024, percentUsed)
		return int(used), int(total), percentUsed, nil

	}
	return 0, 0, 0, fmt.Errorf("unknown error occurred")
}
