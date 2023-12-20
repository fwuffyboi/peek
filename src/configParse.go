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

	// api shit
	Api struct {
		ApiHost string `yaml:"api-host"`
		ApiPort int    `yaml:"api-port"`
	} `yaml:"api"`

	// Action stuff?
	Actions struct {
		SystemShutdown bool `yaml:"system-shutdown"`
		ShutdownPeek   bool `yaml:"shutdown-peek"`
	} `yaml:"actions"`

	// Show stuff?
	Show struct {
		ShowSystemSpecs bool `yaml:"show-system-specs"`
		ShowHostname    bool `yaml:"show-hostname"`
		ShowIP          bool `yaml:"show-ip"`
		ShowCountry     bool `yaml:"show-country"`
		ShowUptime      bool `yaml:"show-uptime"`
		ShowRAM         bool `yaml:"show-ram"`
		ShowCPU         bool `yaml:"show-cpu"`
		ShowDisk        bool `yaml:"show-disk"`
		ShowTemp        bool `yaml:"show-temp"`
		ShowRAID        bool `yaml:"show-raid"`
		ShowGPU         bool `yaml:"show-gpu"`
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
