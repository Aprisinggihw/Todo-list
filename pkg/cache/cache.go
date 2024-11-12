package cache

import (
	"context"
	"todo-list/configs"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitCache(cfg configs.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	return rdb
}

type Cacheable interface {
	Set(key string, value interface{}, duration time.Duration) error
	Get(key string) string
	Delete(key string) error
}

type cacheable struct {
	rdb *redis.Client
}

func NewCacheable(rdb *redis.Client) Cacheable {
	return &cacheable{
		rdb: rdb,
	}
}

func (c *cacheable) Set(key string, value interface{}, duration time.Duration) error {
	err := c.rdb.Set(context.Background(), key, value, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheable) Get(key string) string {
	value, err := c.rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return ""
	}
	return value
}
func (c *cacheable) Delete(key string) error {
    return c.rdb.Del(context.Background(), key).Err()
}