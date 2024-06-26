package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	setLogger()
	params := getCmdLineParams()
	store := NewStore(ctx, params["outputType"])
	cache := NewRedisCache()
	captureCfg := NewPcapCfg(params["device"])

	if err := captureCfg.startPcap(&store, &cache, ctx); err != nil {
		log.Fatalf("could not start pcap %v", err)
	}

}
