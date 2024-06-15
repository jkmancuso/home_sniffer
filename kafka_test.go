package main

import (
	"context"
	"encoding/json"
	"testing"
)

func TestConnect(t *testing.T) {

	kafkaCfg := kafkaConfig{
		topic:     "my-topic",
		partition: 0,
		transport: "tcp",
		host:      "127.0.0.1",
		port:      9094,
	}

	conn, err := kafkaCfg.connectKafka(context.Background())

	if err != nil {
		t.Fatalf("Err: %v\ncould not connect to kafka with params: %+v", err, kafkaCfg)
	}

	store := kafkaStore{
		cfg:  kafkaCfg,
		conn: conn,
	}
	myPacket := packetData{
		Src:    "1.2.3.4",
		Dst:    "5.6.7.8",
		Length: 1,
	}

	t.Run("subtest write", func(t *testing.T) {
		t.Parallel()
		if err := json.NewEncoder(&store).Encode(myPacket); err != nil {
			t.Fatalf("Could not write to store %v", err)
		}

	})

}
