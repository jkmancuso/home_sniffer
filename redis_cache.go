package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const redisEnvfile = "./redis.env"

var errRedisKeyMissing = errors.New("empty key sent")

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
	loadEnv(redisEnvfile)

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cfg := RedisCfg{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	}

	log.Infof("Using redis cache: %+v", cfg)

	return cfg

}

func (r *RedisCfg) newRedisClient() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})

	log.Infof("New redis client created")

	return rdb
}

func NewRedisCache() redisCache {
	cfg := newRedisCfg()
	client := cfg.newRedisClient()

	log.Info("New redis cache returned")

	return redisCache{
		Client: client,
		Cfg:    cfg,
	}

}

func (r redisCache) Get(ctx context.Context, key string) (string, error) {

	if len(key) == 0 {
		log.Error(errRedisKeyMissing.Error())
		return "", errRedisKeyMissing
	}

	result, err := r.Client.Get(ctx, key).Result()

	if err == redis.Nil {
		log.Errorf("key %v not found\n%v", key, err)
		return "", err
	}

	if err != nil {
		//really shouldnt hit this
		errStr := fmt.Sprintf("redis error for key %v", key)
		log.Error(errStr)
		return "", errors.New(errStr)
	}

	log.Debugf("Successfully pulled cache entry %v from redis", result)

	return result, nil

}

func (r redisCache) Set(ctx context.Context, key string, val string) error {
	log.Infof("Adding to redis. key: %v, value: %v", key, val)

	err := r.Client.Set(ctx, key, val, 0).Err()

	if err != nil {
		log.Errorf("Error sending to redis: \nkey: %v\nval: %v\nerr: %v", key, val, err)
		return err
	} else {
		return nil
	}
}
