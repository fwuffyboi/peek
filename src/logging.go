package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func setupLogging() (*os.File, io.Writer, error) {
	// Open or create the log file
	file, err := os.OpenFile("peek.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// Set Logrus to use both the file and stdout as outputs
	multiWriter := io.MultiWriter(file, os.Stdout)
	logrus.SetOutput(multiWriter)

	// TODO: Allow user to choose logging level
	logrus.SetLevel(logrus.InfoLevel)

	// make it json format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return file, multiWriter, nil
}
