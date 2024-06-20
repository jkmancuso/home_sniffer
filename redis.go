package main

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type redisCfg struct {
	Addr     string
	Password string
	DB       int
}

type redisCache struct {
	client *redis.Client
	cfg    redisCfg
}

func newRedisCfg() redisCfg {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cfg := redisCfg{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	}

	log.Debugf("Using redis cache: %+v", cfg)

	return cfg

}

func (r *redisCfg) newRedisClient() *redis.Client {

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
		client: client,
		cfg:    cfg,
	}

}

func (r *redisCache) Get(ctx context.Context, key string) (ipInfo, bool) {
	//result, _ := r.client.Get(ctx,key).Result()

	return ipInfo{}, false
}

func (r *redisCache) Set(ctx context.Context, key string, val string) error {
	log.Debugf("Adding to redis. key: %v, value: %v", key, val)

	err := r.client.JSONSet(ctx, key, "$", val).Err()

	if err != nil {
		log.Errorf("Error sending JSON to redis: \nkey: %v\nval: %v\nerr: %v", key, val, err)
		return err
	} else {
		return nil
	}
}
