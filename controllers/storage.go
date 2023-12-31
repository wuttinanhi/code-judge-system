package controllers

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/redis"
	"github.com/spf13/viper"
)

func GetRedisStorage() fiber.Storage {
	return redis.New(redis.Config{
		Host:      viper.GetString("RATE_LIMIT_HOST"),
		Port:      viper.GetInt("RATE_LIMIT_PORT"),
		Username:  viper.GetString("RATE_LIMIT_USER"),
		Password:  viper.GetString("RATE_LIMIT_PASSWORD"),
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
}

func GetMemoryStorage() fiber.Storage {
	return memory.New(memory.Config{
		GCInterval: 10 * time.Second,
	})
}
