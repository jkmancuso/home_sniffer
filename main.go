package main

import (
	"context"
	"flag"

	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()
	setLogger()

	device := flag.String("device", "br0", "")
	flag.Parse()

	store := NewKafkaStore(ctx)
	//store := newFileStore(ctx)

	cache := NewRedisCache()

	captureCfg := pcapConfig{
		device:  *device,
		snaplen: 1600,
		promisc: true,
		timeout: pcap.BlockForever,
	}

	if err := captureCfg.startPcap(&store, &cache, ctx); err != nil {
		log.Fatalf("could not start pcap %v", err)
	}

}
