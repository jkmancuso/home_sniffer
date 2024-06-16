package main

import (
	"fmt"
	"net"
)

type ip struct {
	ipv4       string
	reverseDNS string
	company    string
}

func newIP(ipAddress string) (ip, error) {

	var reverseDNS, company string
	resultIP := ip{}

	if result := net.ParseIP(ipAddress); result == nil {
		fmt.Printf("%v is not a valid ip", ipAddress)
		return ip{}, nil
	}

	resultIP.ipv4 = ipAddress

	addresses, err := net.LookupAddr(ipAddress)

	//no need to exit
	if err != nil || len(addresses) == 0 {
		fmt.Println(err)
		reverseDNS = ""
	} else {
		reverseDNS = addresses[0]
	}

	resultIP.reverseDNS = reverseDNS
	resultIP.company = company

	return resultIP, nil
}

func getValueFromCache(ipAddress string) string {
	return ipAddress
}
