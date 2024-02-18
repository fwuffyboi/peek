package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func getIP() string {
	url := "https://ipinfo.io/ip"

	log.Info("GeIP: Attempting to get the server's IP address.")
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("GeIP: Unknown error. Err: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("GeIP: Error defering. Err: %s", err)
			return
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		log.Warnf("GeIP: Unexpected status code. Err: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("GeIP: Error reading response body. Err: %s", err)
	}

	config, err := ConfigParser()
	if config.Show.ShowIP == false {
		log.Warn("GeIP: IP address is disabled in the config. Censoring IP address from now.")
		censoredIP := "xxx.xxx.xxx.xxx"
		log.Infof("Server IP address: %s", censoredIP)
	} else {
		log.Infof("Server IP address: %s", string(body))
		return string(body)
	}
}
