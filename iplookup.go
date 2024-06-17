package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

type ipInfo struct {
	ipv4       string
	reverseDNS string
	company    string
}

type registryResult struct {
	Name string `json:"name"`
}

// main function to pack all the lookup info
func GetIPLookupInfo(ipAddress string) (ipInfo, error) {

	info, err := newIPInfo(ipAddress)

	if err != nil {
		fmt.Println(err)
		return ipInfo{}, err
	}

	// if 192 no need to do any lookups
	if info.isLocalLAN() {
		return info, nil
	}

	info.setReverseDNSName()
	info.setCompanyName()

	return info, nil
}

func newIPInfo(ipAddress string) (ipInfo, error) {

	if result := net.ParseIP(ipAddress); result == nil {
		fmt.Printf("%v is not a valid ip", ipAddress)
		return ipInfo{}, nil
	}

	resultIP := ipInfo{
		ipv4: ipAddress,
	}

	return resultIP, nil

}

// basically....is the first octet 192
func (info *ipInfo) isLocalLAN() bool {
	octets := strings.Split(info.ipv4, ".")

	if octets[0] == "192" || octets[0] == "127" {
		return true
	}

	return false

}

func (info *ipInfo) setReverseDNSName() {
	var reverseDNSName string

	addresses, err := net.LookupAddr(info.ipv4)

	//not fatal
	if err != nil || len(addresses) == 0 {
		fmt.Println(err)
		reverseDNSName = ""
	} else {
		reverseDNSName = addresses[0]
	}

	info.reverseDNS = reverseDNSName

}

func (info *ipInfo) setCompanyName() {
	var companyName string

	result, err := queryRegistry(info.ipv4)

	//not fatal
	if err != nil {
		fmt.Println(err)
		companyName = ""
	} else {
		companyName = result.Name
	}

	info.company = companyName
}

func queryRegistry(ipAddress string) (registryResult, error) {
	loadEnv()

	result := registryResult{}

	baseURL := os.Getenv("ARIN_URL")
	url := fmt.Sprintf("%s%s", baseURL, ipAddress)

	fmt.Println("Pulling from arin")
	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return result, err

	}
	fmt.Println("success")

	var buff []byte

	if _, err = response.Body.Read(buff); err != nil {
		fmt.Println(err)
		return result, err
	}

	if err = json.Unmarshal(buff, &result); err != nil {
		fmt.Println(err)
		return result, err
	}

	return result, nil

}
