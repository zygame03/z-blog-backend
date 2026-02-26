package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"my_web/backend/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

func getProfileKey(id int) string {
	return fmt.Sprintf("user:profile:%d", id)
}

func cacheGetProfile(ctx context.Context, rdb *redis.Client, id int) (*Profile, error) {
	key := getProfileKey(id)

	data, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Info(
			"cache miss",
			zap.Int("id", id),
		)
		return nil, ErrCacheMiss
	}

	if err != nil {
		logger.Error(
			"cache get profile failed",
			zap.Int("id", id),
		)
	}

	var data_s Profile
	err = json.Unmarshal([]byte(data), &data_s)
	if err != nil {
		logger.Error(
			"unmarshal failed",
			zap.Int("id", id),
		)
		return nil, err
	}

	return &data_s, nil
}

func cacheSetProfile(ctx context.Context, rdb *redis.Client, id int, data *Profile) error {
	data_s, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	key := getProfileKey(id)
	_, err = rdb.Set(ctx, key, string(data_s), 1*time.Hour).Result()
	if err != nil {
		return err
	}
	return nil
}
