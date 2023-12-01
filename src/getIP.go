package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

func getIP() string {
	url := "https://ipinfo.io/ip"

	log.Info("GeIP: Attempting to get the server's IP address.")
	response, err := http.Get(url)
	if err != nil {
		log.Panicf("GeIP: Unknown error. Err: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("GeIP: Error defering. Err: %s", err)
			return
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		fmt.Printf("GeIP: Unexpected status code. Err: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("GeIP: Error reading response body. Err: %s", err)
	}

	log.Infof("Server IP address: %s", string(body))
	return string(body)
}
