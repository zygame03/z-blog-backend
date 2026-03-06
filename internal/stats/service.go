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
	db  *database
	rdb *cache

	cron *cron.Cron
	ctx  context.Context
}

func NewService(db *gorm.DB, rdb *redis.Client, cron *cron.Cron, ctx context.Context) *Service {
	s := Service{
		db:   newDatabase(db),
		rdb:  NewCache(rdb),
		cron: cron,
		ctx:  ctx,
	}

	cron.AddFunc("@every 10s", s.syncViewUV)
	return &s
}

func (s *Service) syncViewUV() {
	num, err := s.rdb.GetViewUV(s.ctx)
	if err != nil {
		return
	}

	err = s.db.updateViews(num)
	if err != nil {
		return
	}

	err = s.rdb.DelViewUV(s.ctx)
	if err != nil {
		return
	}
}

func (s *Service) RecordUV(ip string) {
	err := s.rdb.AddViewUV(s.ctx, ip)
	if err != nil {
		return
	}
	logger.Info(
		"record",
		zap.String("ip", ip),
	)
}

func (s *Service) GetViews() (int64, int64, error) {
	num, err := s.db.getViews()
	if err != nil {
		return -1, -1, err
	}

	n2, err := s.rdb.GetViewUV(s.ctx)
	if err != nil {
		return num, -1, err
	}

	return num, n2, nil
}
