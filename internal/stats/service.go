package stats

import (
	"context"
	"my_web/backend/internal/global"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db   *repo
	rdb  *cache
	conf func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, conf func() *Config) *Service {
	s := Service{
		db:   newRepo(db),
		rdb:  newCache(rdb, conf),
		conf: conf,
	}

	return &s
}

func (s *Service) RegisterCron(cron *cron.Cron) {
	_, err := cron.AddFunc("@every "+s.conf().SyncInterval.String(), s.syncViewUV)
	if err != nil {
		return
	}
	logger.Info(
		"add func successfully",
		zap.Duration("interval", s.conf().SyncInterval),
		zap.String("func", "SyncInterval"),
	)
}

func (s *Service) syncViewUV() {
	ctx := context.Background()
	num, err := s.rdb.getViewUV(ctx)
	if err != nil {
		logger.Error(
			"get site view uv failed",
			zap.Error(err),
		)
		return
	}

	err = s.db.updateViews(ctx, num)
	if err != nil {
		logger.Error(
			"update site view faied",
			zap.Error(err),
		)
		return
	}

	err = s.rdb.delViewUV(ctx)
	if err != nil {
		logger.Error(
			"delete site view UV failed",
			zap.Error(err),
		)
	}
}

func (s *Service) RecordUV(ctx context.Context, ip string) {
	err := s.rdb.addViewUV(ctx, ip)
	if err != nil {
		logger.Error(
			"record failed",
			zap.String("ip", ip),
		)
		return
	}
	logger.Info(
		"record",
		zap.String("ip", ip),
	)
}

func (s *Service) getViews(ctx context.Context) (int, error) {
	num, err := s.rdb.getView(ctx)
	if err == nil {
		logger.Info(
			"cache get view",
		)
		return num, err
	}
	if err != global.ErrCacheMiss {
		logger.Error(
			"cache get view failed",
			zap.Error(err),
		)
	} else {
		logger.Info(
			"cache miss",
		)
	}

	num, err = s.db.getViews(ctx)
	if err != nil {
		logger.Error(
			"repo get view failed",
			zap.Error(err),
		)
		return -1, err
	}
	err = s.rdb.setView(ctx, num)
	if err != nil {
		logger.Error(
			"cache set view failed",
			zap.Int("view", num),
		)
		return num, err
	}

	return num, nil
}
