package main

import (
	"context"
	"flag"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {

	ctx := context.Background()

	device := flag.String("device", "br0", "")
	flag.Parse()

	store := newKafkaStore(ctx)
	//store := initFileStore(ctx)

	captureCfg := pcapConfig{
		device:  *device,
		snaplen: 1600,
		promisc: true,
		timeout: pcap.BlockForever,
	}

	if err := captureCfg.startPcap(&store); err != nil {
		log.Fatalf("could not start pcap %v", err)
	}

}
