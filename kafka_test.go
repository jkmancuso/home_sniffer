package main

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {

	kafkaCfg := kafkaConfig{
		topic:     "my-topic",
		partition: 0,
		transport: "tcp",
		host:      "localhost",
		port:      9092,
	}

	conn, err := kafkaCfg.connectKafka(context.Background())

	if err != nil {
		t.Fatalf("Err: %v\ncould not connect to kafka with params: %+v", err, kafkaCfg)
	}

	store := kakfkaStore{
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
		err = store.send(&myPacket)

		if err != nil {
			t.Fatalf("Could not write to store %v", err)
		}

	})

}
