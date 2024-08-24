package main

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func TestRedisSetGet(t *testing.T) {

	var ctx = context.Background()

	var tests = []struct {
		testname string
		given    string
		wantStr  string
		wantErr  error
	}{
		{"Empty key", "", "", errRedisKeyMissing},
		{"Key not in redis", "missing_key", "", redis.Nil},
		{"Key is found", "found_key", "123", nil},
	}

	client, mock := redismock.NewClientMock()

	mock.ExpectGet("missing_key").RedisNil()
	mock.ExpectGet("found_key").SetVal("123")

	//pass the mocked client
	r := redisCache{
		Client: client,
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			returnStr, returnErr := r.Get(ctx, test.given)

			if returnStr != test.wantStr || !errors.Is(returnErr, test.wantErr) {
				t.Errorf("Gave key:%v, wanted %v, %v, Got return error: %v",
					returnStr,
					test.wantStr, test.wantErr,
					returnErr)

			}
		})

	}

}

/*
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
