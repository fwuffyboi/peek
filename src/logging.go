package main

import (
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

func setupLogging() (*os.File, io.Writer, error) {
	// grab config settings
	config, err := ConfigParser()
	if err != nil {
		log.Fatal("Error getting config.")
	}

	// Open or create the log file
	var logfileName = config.Logging.LogFile
	usrHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get user's home directory.")
	}

	logfilePath := path.Join(usrHomeDir, ".config/peek", logfileName)
	file, err := os.OpenFile(logfilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// make logrus add this thingy
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			//return frame.Function, fileName
			return "", fileName
		},
	})

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
