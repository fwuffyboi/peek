package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func getUptime() (time.Duration, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		uptimeStr, err := readFileContents("/proc/uptime")
		if err != nil {
			return 0, err
		}

		uptimeSec, err := parseUptime(uptimeStr)
		if err != nil {
			return 0, err
		}

		return time.Duration(uptimeSec) * time.Second, nil

	default:
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func readFileContents(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func parseUptime(uptimeStr string) (float64, error) {
	fields := strings.Fields(uptimeStr)
	if len(fields) < 1 {
		return 0, fmt.Errorf("invalid uptime format")
	}

	uptimeSec, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}

	return uptimeSec, nil
}
