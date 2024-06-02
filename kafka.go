package main

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type kafkaConfig struct {
	topic     string
	partition int
	transport string
	host      string
	port      int
}

type kakfkaStore struct {
	cfg  kafkaConfig
	conn *kafka.Conn
}

func (store kakfkaStore) send(data *packetData) error {
	return nil
}

func (cfg *kafkaConfig) connectKafka(ctx context.Context) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(ctx,
		cfg.transport,
		fmt.Sprintf("%s:%d", cfg.host, cfg.port),
		cfg.topic,
		cfg.partition)

	if err != nil {
		return nil, err
	}

	return conn, err

}
