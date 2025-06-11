package user

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db   *repo
	rdb  *cache
	conf func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, conf func() *Config) *Service {
	return &Service{
		db:   newRepo(db),
		rdb:  newCache(rdb),
		conf: conf,
	}
}

func (s *Service) getProfile(ctx context.Context, id int) (*Profile, error) {
	data, err := s.rdb.getProfile(ctx, id)

	if err == nil {
		logger.Info(
			"get profile from cache",
			zap.Int("id", id),
		)
		return data, nil
	}
	if err != ErrCacheMiss {
		logger.Warn(
			"cache get profile failed",
			zap.Error(err),
		)
	}

	data, err = s.rdb.getProfile(ctx, id)
	if err != nil {
		logger.Error(
			"repo get profile failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		return nil, err
	}

	err = s.rdb.setProfile(ctx, id, data)
	if err != nil {
		logger.Warn(
			"cache set profile failed",
			zap.Int("id", id),
			zap.Error(err),
		)
	}
	return data, nil
}
