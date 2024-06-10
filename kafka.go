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

// Send one msg at a time to kafka, bit inefficient
func (store kafkaStore) sendSingle(data packetData) error {
	kafkaMsg, err := json.Marshal(data)

	if err != nil {
		fmt.Printf("Error marshalling payload, %v", err)
		return err
	}

	_, err = store.conn.Write(kafkaMsg)

	if err != nil {
		fmt.Printf("Error sending payload  to kafka, %v", err)
		return err
	}

	return nil
}

// Send multiple msg at a time to kafka, more efficient
func (store kafkaStore) sendBatch(data []packetData) error {

	var kafkaMsgs []kafka.Message

	for _, payload := range data {

		msgBytes, err := json.Marshal(payload)

		if err != nil {
			fmt.Printf("Error marshalling payload, %v", err)
			return err
		}

		kafkaMsgs = append(kafkaMsgs, kafka.Message{Value: msgBytes})

	}

	_, err := store.conn.WriteMessages(kafkaMsgs...)

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
