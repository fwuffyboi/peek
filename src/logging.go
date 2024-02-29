package main

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func setupLogging() (*os.File, io.Writer, error) {
	// grab config settings
	config, err := ConfigParser()
	if err != nil {
		log.Fatal("Error getting config.")
	}

	// Open or create the log file
	file, err := os.OpenFile("peek.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// Set Logrus to use both the file and stdout as outputs
	multiWriter := io.MultiWriter(file, os.Stdout)
	logrus.SetOutput(multiWriter)

	// This allows the user to choose the logging level
	switch config.Logging.LogLevel {
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERRO":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATA":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		log.Fatalf("Your logging option was invalid. Please choose from INFO, WARN, ERRO and FATA. \nMore information is available at: https://github.com/fwuffyboi/peek#logging-level. Logging option: %s", config.Logging.LogLevel)
	}

	// make it json format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return file, multiWriter, nil
}
