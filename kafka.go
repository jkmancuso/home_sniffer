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

type kafkaStore struct {
	cfg  kafkaConfig
	conn *kafka.Conn
}

// implement io.writer
func (store *kafkaStore) Write(data []byte) (int, error) {
	var packets []packetData
	var kafkaMsgs []kafka.Message
	var err error
	var msgBytes []byte
	var numBytes int

	if err := json.Unmarshal(data, &packets); err != nil {
		fmt.Printf("Error unmarshalling payload, %v", err)
		return 0, err
	}

	for _, payload := range packets {
		msgBytes, err = json.Marshal(payload)

		if err != nil {
			fmt.Printf("Error marshalling payload, %v", err)
			return 0, err
		}

		kafkaMsgs = append(kafkaMsgs, kafka.Message{Value: msgBytes})
	}

	numBytes, err = store.conn.WriteMessages(kafkaMsgs...)

	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	return numBytes, nil
}

// return a kafka connection handle
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
