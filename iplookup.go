package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ipInfo struct {
	Ipv4       string
	ReverseDNS string
	Company    string
}

type registryResult struct {
	Name string `json:"name"`
}

// main function to pack all the lookup info
func GetIPLookupInfo(ipAddress string) (ipInfo, error) {

	info, err := newIPinfo(ipAddress)

	if err != nil {
		log.Errorf("Could not get ip info: %v", err)
		return ipInfo{}, err
	}

	// if 192 no need to do any lookups
	if info.isLocalLAN() {
		return info, nil
	}

	info.setReverseDNSName()
	info.setCompanyName()

	log.Debugf("Returning ip info: %+v", info)

	return info, nil
}

func newIPinfo(ipAddress string) (ipInfo, error) {

	if result := net.ParseIP(ipAddress); result == nil {
		log.Errorf("%v is not a valid ip", ipAddress)
		return ipInfo{}, nil
	}

	resultIP := ipInfo{
		Ipv4: ipAddress,
	}

	return resultIP, nil

}

// basically....is the first octet 192
func (info *ipInfo) isLocalLAN() bool {
	octets := strings.Split(info.Ipv4, ".")

	if octets[0] == "192" || octets[0] == "127" {
		return true
	}

	return false

}

func (info *ipInfo) setReverseDNSName() {
	var reverseDNSName string

	addresses, err := net.LookupAddr(info.Ipv4)

	//not fatal
	if err != nil || len(addresses) == 0 {
		log.Errorf("Error looking up %v\n%v", info.Ipv4, err)
		reverseDNSName = ""
	} else {
		reverseDNSName = addresses[0]
	}

	log.Debugf("Got reverseDNS Name: %v", reverseDNSName)

	info.ReverseDNS = reverseDNSName

}

// get name of company from registrar
func (info *ipInfo) setCompanyName() {
	var companyName string

	result, err := queryRegistry(info.Ipv4)

	//not fatal
	if err != nil {
		log.Errorf("Error query registry %v", err)
		companyName = ""
	} else {
		companyName = result.Name
	}

	log.Debugf("Got company name: %v", companyName)

	info.Company = companyName
}

// query arin or equivalent registrar for info
func queryRegistry(ipAddress string) (registryResult, error) {
	loadEnv()

	result := registryResult{}

	baseURL := os.Getenv("ARIN_URL")
	url := fmt.Sprintf("%s%s", baseURL, ipAddress)

	log.Debugf("Query registry %s:", url)

	response, err := http.Get(url)

	if err != nil {
		log.Error(err)
		return result, err
	}

	log.Debug("Success!")

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Errorf("Failed to read response: %v", err)
		return result, err
	}

	//log.Debugf("Got response: %s", body)

	defer response.Body.Close()

	if err = json.Unmarshal(body, &result); err != nil {
		log.Errorf("Unable to unmarshall: %v\n%s", err, body)
		return result, err
	}

	log.Debugf("Registry result: %+v", result)

	return result, nil

}
