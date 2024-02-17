package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type T struct {
	Auth struct {
		AuthRequired bool   `yaml:"auth-required"`
		AuthLevel    int    `yaml:"auth-level"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
	} `yaml:"auth"`

	// log settings
	Logging struct {
		LogLevel string `yaml:"log-level"`
		LogFile  string `yaml:"log-file"`
	} `yaml:"logging"`

	// api shit
	Api struct {
		ApiHost   string `yaml:"api-host"`
		ApiPort   int    `yaml:"api-port"`
		RateLimit int    `yaml:"rate-limit"`
	} `yaml:"api"`

	// Action stuff?
	Actions struct {
		SystemShutdown bool `yaml:"system-shutdown"`
		ShutdownPeek   bool `yaml:"shutdown-peek"`
	} `yaml:"actions"`

	// Show stuff?
	Show struct {
		ShowSystemSpecs   bool `yaml:"show-sys-specs"`
		ShowHostname      bool `yaml:"show-hostname"`
		ShowIP            bool `yaml:"show-ip"`
		ShowServerCountry bool `yaml:"show-server-country"`
		ShowClientCountry bool `yaml:"show-client-country"`
		ShowUptime        bool `yaml:"show-uptime"`
		ShowRAM           bool `yaml:"show-ram"`
		ShowCPU           bool `yaml:"show-cpu"`
		ShowDisk          bool `yaml:"show-disk"`
		ShowCPUTemp       bool `yaml:"show-cpu-temp"`
		ShowCPUUsage      bool `yaml:"show-cpu-use"`
		ShowRAID          bool `yaml:"show-raid"`
		ShowGPU           bool `yaml:"show-gpu"`
		ShowLogsAPI       bool `yaml:"show-logs-api"`
		ShowErrors        bool `yaml:"show-errors"`
	} `yaml:"show"`
}

func ConfigParser() (*T, error) {
	yamlFile, err := os.ReadFile("peek.config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	var t T
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	return &t, nil

}
