package database

import (
	"backend/internal/utils"
	"context"
	"os"

	"github.com/gofiber/storage/redis/v3"
)


func NewRedis() *redis.Storage{
	redisStorage := redis.New(redis.Config{
		URL: os.Getenv("REDIS_URL"),
		Port: 6379,
		Password: os.Getenv("REDIS_PASSWORD"),
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10,
	})

	return redisStorage
}

type RedisLocker struct {
	Redis *redis.Storage
}

func (r *RedisLocker) Lock(key string) error {
    err := r.Redis.SetWithContext(context.Background(), "lock:"+key, make([]byte, 1), utils.ThirtySec)
    if err != nil {
        return err
    }
    return nil
}

func (r *RedisLocker) Unlock(key string) error {
    return r.Redis.DeleteWithContext(context.Background(), "lock:"+key)
}