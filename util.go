package main

import (
	"regexp"
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
