package main

import (
	"context"
	"encoding/json"
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
	kafkaMsg, _ := json.Marshal(data)

	_, err := store.conn.Write(kafkaMsg)

	if err != nil {
		fmt.Printf("Error sending payload  to kafka, %v", err)
		return err
	}

	return nil
}

func (cfg *kafkaConfig) connectKafka(ctx context.Context) (*kafka.Conn, error) {

	fmt.Printf("Connecting to kafka\n%+v", cfg)
	conn, err := kafka.DialLeader(ctx,
		cfg.transport,
		fmt.Sprintf("%s:%d", cfg.host, cfg.port),
		cfg.topic,
		cfg.partition)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Success")

	return conn, err

}
