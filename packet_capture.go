package main

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type pcapConfig struct {
	device  string
	snaplen int32
	promisc bool
	timeout time.Duration
}

// Start new packet capture
func (cfg *pcapConfig) startPcap() error {
	fmt.Printf("Starting packet cap on device %v\n", cfg.device)

	handle, err := cfg.newPcapHandle()

	if err != nil {
		return err
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		fmt.Printf("%v", packet)
	}

	return nil

}

// Generate a new pcap handle given a config
func (cfg *pcapConfig) newPcapHandle() (*pcap.Handle, error) {

	if err := cfg.validateInterfaces(); err != nil {
		return nil, err
	}

	handle, err := pcap.OpenLive(
		cfg.device,
		cfg.snaplen,
		cfg.promisc,
		cfg.timeout,
	)

	if err != nil {
		return nil, err
	}

	return handle, err

}

// Helper function to check that the selected device (ie br0 or eth0) is valid on your machine
func (cfg *pcapConfig) validateInterfaces() error {
	fmt.Println("Confirming valid devices")

	selectedDevice := cfg.device

	ifs, err := pcap.FindAllDevs()

	if err != nil {
		return err
	}

	for _, ethInterface := range ifs {
		if ethInterface.Name == selectedDevice {
			fmt.Printf("Found device %v, checking IP...\n", selectedDevice)

			ips := ethInterface.Addresses

			if len(ips) > 0 {
				fmt.Printf("IP: %v", ips)
				return nil
			}

		}
	}

	return fmt.Errorf("Interface is not valid %w", selectedDevice)

}
