package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// this is used to represent 1 connection essentially
type ipInfo struct {
	Ipv4 string
	DNS  string
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

	info := ipInfo{
		Ipv4: ipAddress,
	}

	log.Debugf("Getting ip info for %v", ipAddress)

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// if 192 no need to do any lookups
	if info.isLocalLAN() {
		return info, nil
	}

	DNSname, found := cache.Get(ctx, ipAddress)

	info.DNS = DNSname

	if found { // if its in the cache you're good. just return info
		log.Debugf("Found IP %v in cache!", info.Ipv4)
	} else { // if its not in cache...
		log.Debugf("Did not find IP %v in cache", info.Ipv4)
	}

	log.Debugf("Returning ip info: %+v", info)

	// then return the info struct
	return info, nil
}

func NewIPinfo(ipAddress string, cache Cache, ctx context.Context) (ipInfo, error) {

	if result := net.ParseIP(ipAddress); result == nil {
		log.Errorf("%v is not a valid ip", ipAddress)
		return ipInfo{}, errors.New("invalid ip")
	}

	resultInfo, _ := GetIPLookupInfo(ipAddress, cache, ctx)

	log.Debugf("Created new ipinfo struct %v", resultInfo)

	return resultInfo, nil

}

// basically....is the first octet 192
func (info *ipInfo) isLocalLAN() bool {
	octets := strings.Split(info.Ipv4, ".")

	if octets[0] == "192" || octets[0] == "127" {
		return true
	}

	return false

}
