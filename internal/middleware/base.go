package middleware

import (
	"my_web/backend/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constants.DATABASE, db)
		ctx.Next()
	}
}

func WithRedis(redis *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constants.REDIS, redis)
		ctx.Next()
	}
}
