package user

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	db   *repo
	rdb  *cache
	conf func()
}

func NewService(db *gorm.DB, rdb *redis.Client, conf func()) *service {
	return &service{
		db:   newRepo(db),
		rdb:  newCache(rdb),
		conf: conf,
	}
}

func (s *service) getProfile(ctx context.Context, id int) (*Profile, error) {
	data, err := s.rdb.getProfile(ctx, id)

	if err == nil {
		logger.Info(
			"get profile from cache",
			zap.Int("id", id),
		)
		return data, nil
	}

	data, err = s.rdb.getProfile(ctx, id)
	if err != nil {
		logger.Error(
			"repo get profile failed",
			zap.Int("id", id),
		)
		return nil, err
	}

	s.rdb.setProfile(ctx, id, data)
	return data, nil
}
