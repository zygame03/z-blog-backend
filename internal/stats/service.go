package stats

import (
	"context"
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
		rdb:  newCache(rdb),
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
		"添加定时任务成功",
		zap.Int("间隔", int(s.conf().SyncInterval)),
		zap.String("任务描述", "网站浏览数同步"),
	)
}

func (s *Service) syncViewUV() {
	ctx := context.Background()
	num, err := s.rdb.getViewUV(ctx)
	if err != nil {
		logger.Error(
			"get site view uv failed",
		)
		return
	}

	err = s.db.updateViews(num)
	if err != nil {
		logger.Error(
			"update site view failed",
		)
		return
	}

	err = s.rdb.delViewUV(ctx)
	if err != nil {
		logger.Error(
			"delete site view UV failed",
		)
	}
}

func (s *Service) RecordUV(ctx context.Context, ip string) {
	err := s.rdb.addViewUV(ctx, ip)
	if err != nil {
		return
	}
	logger.Info(
		"record",
		zap.String("ip", ip),
	)
}

func (s *Service) getViews() (int64, error) {
	num, err := s.db.getViews()
	if err != nil {
		return -1, err
	}
	return num, nil
}
