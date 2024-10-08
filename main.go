package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	setLogger()
	params := getCmdLineParams()
	store, err := NewStore(ctx, *params["outputType"])

	if err != nil {
		log.Fatal("could not connect to output store!", err)
	}

	cache := NewCache(*params["cacheType"])
	captureCfg := NewPcapCfg(params)

	if err := captureCfg.startPcap(ctx, store, cache); err != nil {
		log.Fatalf("could not start pcap %v", err)
	}

}
