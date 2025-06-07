package home

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db  *database
	rdb *cache

	getCfg func()
}

func NewService(db *gorm.DB, rdb *redis.Client, getCfg func()) *Service {
	service := &Service{
		getCfg: getCfg,
	}

	service.db = newDatabase(db)
	service.rdb = newCache(rdb, service.getCfg)

	return service
}
