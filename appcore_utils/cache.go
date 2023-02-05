package appcore_utils

import (
	"context"

	"github.com/go-redis/redis/v8"
)

//redis

var ctx = context.Background()

func NewCache(configs *Configurations) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.RedisUrl,
		Password: configs.RedisPass,
		DB:       0,
	})

	status := rdb.Ping(ctx)
	if status.Err() != nil {
		panic("cannot connect redis database >> " + status.Err().Error())
	}

	return rdb
}
