package main

/*
import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestRedisSetGet(t *testing.T) {
	loadEnv()
	client := NewRedisCache()

	log.Printf("Using redis cfg %v", client.Cfg)

	key := "1.2.3.4"
	sentipInfo := ipInfo{
		Ipv4:       key,
		ReverseDNS: "mycompany.com",
		Company:    "JOHNS AWESOME COMPANY",
	}

	val, _ := json.Marshal(sentipInfo)

	err := client.Set(context.Background(), key, string(val))

	if err != nil {
		t.Errorf("Unable to set redis key %s\n%v", val, err)
	}

	returnedipInfo, found := client.Get(context.Background(), key)

	if !found {
		t.Errorf("Could not find redis key %s", val)
	}

	if returnedipInfo != sentipInfo {
		t.Errorf("Sent: %v\nGot: %v", sentipInfo, returnedipInfo)
	}

	t.Cleanup(func() {
	})

}
*/
