package main

import (
	"flag"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {

	device := flag.String("device", "br0", "")
	flag.Parse()

	cfg := pcapConfig{
		device:  *device,
		snaplen: 1600,
		promisc: true,
		timeout: pcap.BlockForever,
	}

	if err := cfg.startPcap(); err != nil {
		log.Fatalf("could not start pcap %v", err)
	}

}
