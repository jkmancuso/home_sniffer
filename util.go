package main

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func parseIPs(payload string) (string, string) {
	r, _ := regexp.Compile(`SrcIP:(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}).*DstIP:(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	matches := r.FindStringSubmatch(payload)

	if len(matches) != 3 {
		return "", ""
	}

	return matches[1], matches[2]

}

func parseSize(payload string) string {
	r, _ := regexp.Compile(` Length:(\d+)`)
	matches := r.FindStringSubmatch(payload)

	if len(matches) != 2 {
		return ""
	}

	return matches[1]

}

func loadEnv() {
	godotenv.Load()
}

func setLogger() {
	loadEnv()

	level, err := log.ParseLevel(os.Getenv("LOGLEVEL"))

	if err != nil {
		log.Panic("Unable to recoghnize logging")
	}

	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	//log.SetReportCaller(true)

	log.Printf("log level set to: %v", level.String())

}
