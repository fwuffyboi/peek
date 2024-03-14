package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
		ApiHost       string `yaml:"api-host"`
		ApiPort       int    `yaml:"api-port"`
		RateLimit     int    `yaml:"rate-limit"`
		ShutdownDelay int    `yaml:"shutdown-delay"`
	} `yaml:"api"`

	// Action stuff?
	Actions struct {
		SystemShutdown bool `yaml:"shutdown-system"`
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
		ShowTimezone      bool `yaml:"show-timezone"`
	} `yaml:"show"`

	// integrations
	Integrations struct {
		Telegram struct {
			Enabled          bool   `yaml:"telegram-enabled"`
			TelegramBotToken string `yaml:"telegram-bot-token"`
			TelegramChatID   string `yaml:"telegram-chat-id"`
			TelegramMessage  string `yaml:"telegram-message"`
		} `yaml:"telegram"`
	}
}

// ConfigParser func for parsing the config file and returning all values
func ConfigParser() (*T, error) {
	// grab home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home dir: %v", err)
	}

	// join home dir with filepath to config
	cfFilePath := filepath.Join(homeDir, ".config/peek/peek.config.yaml")

	// read
	if _, err := os.Stat(cfFilePath); errors.Is(err, os.ErrNotExist) {
		// config does not exist or cant be accessed.
		log.Warn("peek.config.yaml not found, creating...")
		err = makeConfig()
		if err != nil {
			log.Fatalf("Error creating default config file: %v", err)
			return nil, fmt.Errorf("error creating config: %v", err)
		}
	}

	yamlFile, err := os.ReadFile(cfFilePath)
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
		return nil, fmt.Errorf("error reading YAML file: %v", err)

	}
	// unmarshal and make into vars
	var t T
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	// give back to whatever called it
	return &t, nil

}

func makeConfig() error {
	// make directory .config/peek
	log.Warn("makeConfig() has been called")
	homeDir, err := os.UserHomeDir()
	log.Info("User home dir: ", homeDir)
	if err != nil {
		log.Fatalf("error getting home dir: %v", err)
		return fmt.Errorf("error getting home dir: %v", err)
	}

	dirToMakePath := filepath.Join(homeDir, ".config/peek")
	log.Infof("Directory to make: %s", dirToMakePath)
	err = os.MkdirAll(dirToMakePath, 0755)
	if err != nil {
		log.Fatalf("error making directory: %v", err)
		return fmt.Errorf("error making directory: %v", err)
	}

	// download from GitHub
	log.Infof("Downloading peek.config.yaml from GitHub...")
	fileURL := "https://raw.githubusercontent.com/fwuffyboi/peek/master/peek.config.yaml.DOWNLOAD"
	destFilePath := filepath.Join(homeDir, ".config/peek/peek.config.yaml")

	log.Info("Creating destination file...")
	file, err := os.Create(destFilePath)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close() // todo: error handle this

	resp, err := http.Get(fileURL)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close() // todo: error handle this x2

	// Check if the response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to download file, http status not 200: ", resp.Status)
		return fmt.Errorf("failed to download file: %v", resp.Status)
	}

	// write to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Error copying file content: %v", err)
		return fmt.Errorf("error copying file content: %v", err)
	}

	// return nil if no errors
	log.Infof("peek.config.yaml has been successfully downloaded to %s", destFilePath)

	return nil
}
