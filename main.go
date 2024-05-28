package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/gopacket/pcap"
)

func main() {

	device := flag.String("device", "br0", "")
	flag.Parse()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg := pcapConfig{
		device:  *device,
		snaplen: 1600,
		promisc: true,
		timeout: pcap.BlockForever,
	}

	if err := cfg.startPcap(); err != nil {
		log.Fatal("Could not start pcap %v", err)
	}

}
