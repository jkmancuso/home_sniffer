package main

import (
	"context"
	"testing"
)

func TestConnectAndSend(t *testing.T) {

	storeType := "kafka"
	kafkaStore, err := NewStore(context.Background(), storeType)

	if err != nil {
		t.Fatalf("Err: %v\ncould not connect to %v", err, storeType)
	}

	sent := []string{"1.2.3.4"}

	t.Run("subtest write", func(t *testing.T) {
		t.Parallel()
		if err := kafkaStore.Send(sent); err != nil {
			t.Fatalf("Could not write to store %v", err)
		}
	})

}
