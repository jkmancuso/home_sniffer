package main

import (
	"flag"
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

func loadEnv(paths ...string) {
	if len(paths) == 0 {
		godotenv.Load()
	} else {
		for _, path := range paths {
			godotenv.Load(path)
		}
	}

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

func getCmdLineParams() map[string]string {
	params := make(map[string]string)

	params["device"] = *flag.String("device", "wlan0", "")
	params["outputType"] = *flag.String("output", "kafka", "")
	params["cacheType"] = *flag.String("cache", "redis", "")
	flag.Parse()

	log.Debugf("Get flags: %+v", params)
	return params
}
