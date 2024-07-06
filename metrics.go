package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

type metrics struct {
	cache   *prometheus.GaugeVec
	traffic *prometheus.GaugeVec
}

type gaugeVec struct {
	opts   prometheus.GaugeOpts
	labels []string
}

func NewMetrics() *metrics {

	log.Info("Creating new prometheus metrics")

	cacheGaugeVec := gaugeVec{
		opts:   prometheus.GaugeOpts{Name: "cache_get"},
		labels: []string{"hit"},
	}

	trafficGaugeVec := gaugeVec{
		opts:   prometheus.GaugeOpts{Name: "traffic"},
		labels: []string{"src", "dst"},
	}

	m := &metrics{
		cache:   promauto.NewGaugeVec(cacheGaugeVec.opts, cacheGaugeVec.labels),
		traffic: promauto.NewGaugeVec(trafficGaugeVec.opts, trafficGaugeVec.labels),
	}

	return m
}
