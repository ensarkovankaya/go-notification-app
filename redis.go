package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Redis *redis.Client

func closeRedis() {
	if err := Redis.Close(); err != nil {
		zap.L().Error("failed to close redis connection", zap.Error(err))
	}
}

func init() {
	options, err := redis.ParseURL(Cnf.RedisURI)
	if err != nil {
		panic(fmt.Errorf("failed to parse redis uri '%v': %w", Cnf.RedisURI, err))
	}

	Redis = redis.NewClient(options)
}
