package infra

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Protocol int    `mapstructure:"protocol"`
}

// InitCache 初始化 Redis 连接
func InitCache(cfg *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	log.Println("Redis 初始化成功")
	return rdb, nil
}
