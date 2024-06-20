package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

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

func (info *ipInfo) String() string {
	infoBytes, err := json.Marshal(info)

	if err != nil {
		return ""
	}

	return string(infoBytes)

}

// main function to pack all the lookup info
func GetIPLookupInfo(ipAddress string, cache Cache, ctx context.Context) (ipInfo, error) {

	var reverseDNSName, companyName string

	log.Debugf("Getting ip info for %v", ipAddress)

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	info, err := newIPinfo(ipAddress)

	if err != nil {
		log.Errorf("Could not get ip info: %v", err)
		return ipInfo{}, err
	}

	// if 192 no need to do any lookups
	if info.isLocalLAN() {
		return info, nil
	}

	//info, found := cache.Get(ctx, info.Ipv4)

	found := false

	if found { // if its in the cache you're good. just return info
		log.Debugf("Found IP %v in cache!", info.Ipv4)
	} else { // if its not in cache...
		log.Debugf("Did not find IP %v in cache", info.Ipv4)

		// do the manual lookup*
		reverseDNSName = info.lookupReverseDNSName()
		companyName = info.lookupCompanyName()

		// update the info struct
		info.setCompanyName(companyName)
		info.setReverseDNSName(reverseDNSName)

		// update the cache
		_ = cache.Set(ctx, info.Ipv4, info.String())
	}

	log.Debugf("Returning ip info: %+v", info)

	// then return the info struct
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

	log.Debugf("Created new ipinfo struct %v", resultIP)

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

func (info *ipInfo) lookupReverseDNSName() string {
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

	return reverseDNSName

}

// get name of company from registrar
func (info *ipInfo) lookupCompanyName() string {
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

	return companyName
}

func (info *ipInfo) setCompanyName(name string) {
	info.Company = name
}

func (info *ipInfo) setReverseDNSName(name string) {
	info.ReverseDNS = name
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
