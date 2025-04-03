package user

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	db  *gorm.DB
	rdb *redis.Client
}

func newService(db *gorm.DB, rdb *redis.Client) *service {
	return &service{
		db:  db,
		rdb: rdb,
	}
}

func (s *service) getProfile(ctx context.Context, id int) (*Profile, error) {
	data, err := cacheGetProfile(ctx, s.rdb, id)

	if err == nil {
		logger.Info(
			"get profile from cache",
			zap.Int("id", id),
		)
		return data, nil
	}

	data, err = repoGetProfile(s.db, id)
	if err != nil {
		logger.Error(
			"repo get profile failed",
			zap.Int("id", id),
		)
		return nil, err
	}

	cacheSetProfile(ctx, s.rdb, id, data)
	return data, nil
}
