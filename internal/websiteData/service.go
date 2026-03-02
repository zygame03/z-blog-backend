package websiteData

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type service struct {
	db  *database
	rdb *cache

	getCfg func() *WebsiteDataConfig
}

func newService(db *gorm.DB, rdb *redis.Client, cfg func() *WebsiteDataConfig) *service {
	s := service{
		getCfg: cfg,
	}

	s.db = NewDatabase(db)
	s.rdb = NewCache(rdb, cfg)

	return &s
}

func (s *service) getIntro(ctx context.Context) (string, error) {
	data, err := s.rdb.GetIntro(ctx)
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

		s.rdb.SetIntro(ctx, data)
		return data, err
	}

	return s.db.getIntro()
}

func (s *service) getAnnouncement(ctx context.Context) ([]string, error) {
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
