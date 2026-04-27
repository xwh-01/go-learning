package config

import (
	"context"
	"exchangeapp/global"
	"exchangeapp/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

func initRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		DB:       AppConfig.Redis.DB,
		Password: AppConfig.Redis.Password,
		PoolSize: AppConfig.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		utils.Fatal("Failed to connect to Redis, got error: %v", err)
	}

	global.RedisDB = RedisClient
	utils.Info("Redis连接成功")
}
