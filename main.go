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

	kafkaCfg := kafkaConfig{
		topic:     "my-topic",
		partition: 0,
		transport: "tcp",
		host:      "localhost",
		port:      9094,
	}

	conn, err := kafkaCfg.connectKafka(ctx)

	if err != nil {
		log.Fatalf("Err: %v\ncould not connect to kafka with params: %+v", err, kafkaCfg)
	}

	store := kafkaStore{
		cfg:  kafkaCfg,
		conn: conn,
	}

	captureCfg := pcapConfig{
		device:  *device,
		snaplen: 1600,
		promisc: true,
		timeout: pcap.BlockForever,
	}

	if err := captureCfg.startPcap(store); err != nil {
		log.Fatalf("could not start pcap %v", err)
	}

}
