package site

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type announcementBO struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

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

func (s *Service) getAnnouncement(ctx context.Context) ([]*announcementBO, error) {
	data, err := s.rdb.getAnnouncement(ctx)
	if err == nil {
		logger.Info(
			"get announcement from cache",
		)

		return data, err
	}

	if err == ErrCacheMiss {
		logger.Info(
			"cache miss for announcement",
		)
		data, err = s.db.getAnnouncement(ctx)
		if err != nil {
			logger.Error(
				"repo get announcement failed",
				zap.Error(err),
			)
			return data, err
		}

		s.rdb.setAnnouncement(ctx, data)
		return data, err
	}

	return s.db.getAnnouncement(ctx)
}
