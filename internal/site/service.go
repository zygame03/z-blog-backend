package site

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db  *repo
	rdb *cache

	getCfg func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, cfg func() *Config) *Service {
	s := Service{
		getCfg: cfg,
	}

	s.db = newRepo(db)
	s.rdb = newCache(rdb, cfg)

	return &s
}

func (s *Service) getIntro(ctx context.Context) (string, error) {
	data, err := s.rdb.getIntro(ctx)
	if err == nil {
		logger.Info(
			"get intro from cache",
		)
		return data, err
	}

	if err == ErrCacheMiss {
		logger.Info(
			"cache miss for intro",
		)
		data, err = s.db.getIntro()
		if err != nil {
			return data, err
		}

		s.rdb.setIntro(ctx, data)
		return data, err
	}

	return s.db.getIntro()
}

func (s *Service) getAnnouncement(ctx context.Context) ([]string, error) {
	data, err := s.rdb.getAnnouncement(ctx)
	if err == nil {
		logger.Info(
			"get intro from cache",
		)

		return data, err
	}

	if err == ErrCacheMiss {
		logger.Info(
			"cache miss for intro",
		)
		data, err = s.db.getAnnouncement()
		if err != nil {
			return data, err
		}

		s.rdb.setAnnouncement(ctx, data)
		return data, err
	}

	return s.db.getAnnouncement()
}
