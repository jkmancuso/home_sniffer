package main

import (
	"fmt"
	"strconv"
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

type packetData struct {
	Src    string
	Dst    string
	Length int
}

// Start new packet capture
func (cfg *pcapConfig) startPcap(store packetStore) error {
	fmt.Printf("Starting packet cap on device %v\n", cfg.device)

	handle, err := cfg.newPcapHandle()

	if err != nil {
		return err
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	var netLayer, transportLayer, srcIP, dstIP, size string
	var i int64

	batchSize := 10
	packetBatch := []packetData{}

	for packet := range packetSource.Packets() {

		//experiencing random SEGFAULT when grabbing netflow
		//https://pkg.go.dev/github.com/google/gopacket#NetworkLayer
		//so have to parse as string :(
		netLayer = fmt.Sprintf("%+v", packet.NetworkLayer())
		transportLayer = fmt.Sprintf("%+v", packet.TransportLayer())

		srcIP, dstIP = parseIPs(netLayer)
		size = parseSize(transportLayer)

		fmt.Printf("src:%v,dst:%v,size:%v\n", srcIP, dstIP, size)

		//some packets have no payload such as ACKs, just move to next iteration
		if len(size) == 0 {
			fmt.Println("Dropping")
			continue
		}

		i += 1

		sizeInt, _ := strconv.Atoi(size)

		pack := packetData{
			Src:    srcIP,
			Dst:    dstIP,
			Length: sizeInt,
		}

		packetBatch = append(packetBatch, pack)

		if i%int64(batchSize) == 0 {

			if err := store.sendBatch(packetBatch); err != nil {
				fmt.Println(err)
			}

			packetBatch = packetBatch[:0]
		}

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

	return fmt.Errorf("interface is not valid %v", selectedDevice)

}
