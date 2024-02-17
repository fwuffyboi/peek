package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"strings"
)

func getMemoryUsage() (total uint64, free uint64, used uint64, usedPercent int, err error) {

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
