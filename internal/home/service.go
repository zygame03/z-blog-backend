package home

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	getCfg func()
}

func NewService(db *gorm.DB, rdb *redis.Client, getCfg func()) *Service {
	service := &Service{
		getCfg: getCfg,
	}

	return service
}
