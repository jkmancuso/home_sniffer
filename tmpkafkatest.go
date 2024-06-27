package main

/*
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

	testipInfo := ipInfo{
		Ipv4:       "1.2.3.4",
		Company:    "mycompany",
		ReverseDNS: "reverse.mycompany.com",
	}

	myPacket := packetData{
		Src:    testipInfo,
		Dst:    testipInfo,
		Length: 1,
	}

	myPackets := []packetData{myPacket}

	t.Run("subtest write", func(t *testing.T) {
		t.Parallel()
		if err := json.NewEncoder(&store).Encode(myPackets); err != nil {
			t.Fatalf("Could not write to store %v", err)
		}

	})

}
*/
