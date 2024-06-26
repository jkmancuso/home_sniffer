package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisCfg struct {
	Addr     string
	Password string
	DB       int
}

type redisCache struct {
	Client *redis.Client
	Cfg    RedisCfg
}

func newRedisCfg() RedisCfg {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cfg := RedisCfg{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	}

	log.Debugf("Using redis cache: %+v", cfg)

	return cfg

}

func (r *RedisCfg) newRedisClient() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})

	log.Debug("New redis client created")

	return rdb
}

func NewRedisCache() redisCache {
	cfg := newRedisCfg()
	client := cfg.newRedisClient()

	log.Debug("New redis cache returned")

	return redisCache{
		Client: client,
		Cfg:    cfg,
	}

}

func (r redisCache) Get(ctx context.Context, key string) (ipInfo, bool) {

	resultipInfo := ipInfo{
		Ipv4: key,
	}

	result, err := r.Client.JSONGet(ctx, key).Result()

	if err != nil {
		log.Errorf("Got error from redis getting key %v\n%v", key, err)
		return resultipInfo, false
	}

	err = json.Unmarshal([]byte(result), &resultipInfo)

	if err != nil {
		log.Errorf("Cannot unmarshall string from redis: %v\n%v", key, err)
		return resultipInfo, false
	}

	log.Debugf("Successfully pulled cache entry %v from redis", resultipInfo)

	return resultipInfo, true

}

func (r redisCache) Set(ctx context.Context, key string, val string) error {
	log.Debugf("Adding to redis. key: %v, value: %v", key, val)

	err := r.Client.JSONSet(ctx, key, "$", val).Err()

	if err != nil {
		log.Errorf("Error sending JSON to redis: \nkey: %v\nval: %v\nerr: %v", key, val, err)
		return err
	} else {
		return nil
	}
}
