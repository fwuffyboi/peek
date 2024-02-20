package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oschwald/maxminddb-golang"
	log "github.com/sirupsen/logrus"
)

func downloadIPDB() error { // todo IMPROVE THIS, 1ST PRIORITY
	// make directory .config/peek
	log.Warn("downloadIPDB() has been called")
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
	log.Infof("Downloading dbip-country-lite-2023-11.mmdb from GitHub...")
	fileURL := "https://raw.githubusercontent.com/fwuffyboi/peek/master/src/assets/dbip-country-lite-2023-11.mmdb"
	destFilePath := filepath.Join(homeDir, ".config/peek/dbip-country-lite-2023-11.mmdb")

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
	log.Infof("The IP database has been successfully downloaded to %s", destFilePath)

	return nil

}

func countryFromIP(ipAddress string) string {
	// put path together todo
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("CFIP: Error getting home dir. Err: %s", err)
		return "Unknown"
	}
	dbPath := filepath.Join(homeDir, ".config/peek/dbip-country-lite-2023-11.mmdb")

	// try to access db
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		log.Errorf("CFIP: Err: %s", err)
		// assume not found, download it todo
		err = downloadIPDB()
		if err != nil {
			log.Fatal("Unable to download the IPDB database. Error: ", err)
		}
	}
	defer func(db *maxminddb.Reader) {
		err := db.Close()
		if err != nil {
			log.Errorf("CFIP: Error defering. err: %s", err)
		}
	}(db)
	// try again
	db, err = maxminddb.Open(dbPath)
	if err != nil {
		log.Errorf("CFIP: Err: %s", err)
	}
	defer func(db *maxminddb.Reader) {
		err := db.Close()
		if err != nil {
			log.Errorf("CFIP: Error defering. err: %s", err)
		}
	}(db)

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		log.Error("CFIP: No IP address")
	}

	err = db.Lookup(ip, &record)
	if err != nil {
		log.Errorf("CFIP: Error occured while looking IP up in database. Err: %s", err)
	}

	if record.Country.ISOCode == "" {
		if ipAddress == "127.0.0.1" {
			log.Infof("CFIP: IP: %s, Country: %s", ipAddress, "Localhost")
			return "Localhost"
		} else {
			log.Warnf("CFIP: IP: %s, Country: %s", ipAddress, "Unknown")
			return "Unknown"
		}
	} else {
		log.Infof("CFIP: IP: %s, Country: %s", ipAddress, record.Country.ISOCode)
		return record.Country.ISOCode
	}
}
