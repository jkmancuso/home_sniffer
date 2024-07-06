package main

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	setLogger()
	params := getCmdLineParams()
	store, err := NewStore(ctx, params["outputType"])

	if err != nil {
		log.Fatal("could not connect to output store!", err)
	}

	cache := NewCache(params["cacheType"])
	captureCfg := NewPcapCfg(params)

	m := NewMetrics()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	if err := captureCfg.startPcap(ctx, store, cache, m); err != nil {
		log.Fatalf("could not start pcap %v", err)
	}

}
