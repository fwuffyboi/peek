package main

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
	log "github.com/sirupsen/logrus"
)

func countryFromIP(ipAddress string) string {
	db, err := maxminddb.Open("./src/assets/dbip-country-lite-2023-11.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *maxminddb.Reader) {
		err := db.Close()
		if err != nil {
			log.Fatalf("CFIP: Error defering. err: %s", err)
		}
	}(db)

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		log.Fatal("CFIP: No IP address")
	}

	err = db.Lookup(ip, &record)
	if err != nil {
		log.Fatalf("CFIP: Error occured while looking IP up in database. Err: %s", err)
	}

	if record.Country.ISOCode == "" {
		if ipAddress == "127.0.0.1" {
			log.Infof("CFIP: IP: %s, Country: %s", ipAddress, "Localhost")
			return "Localhost"
		} else {
			log.Warnf("CFIP: IP: %s, Country: %s", ipAddress, record.Country.ISOCode)
			return "Unknown"
		}
	} else {
		log.Infof("CFIP: IP: %s, Country: %s", ipAddress, record.Country.ISOCode)
		return record.Country.ISOCode
	}
}
