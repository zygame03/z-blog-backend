package data

import (
	"my_web/backend/internal/httpserver"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	httpserver.BaseHandler
	service *service
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api")
	r.POST("/config", middleware.JWTAuth(), h.changeConfig)
	r.GET("/data/intro", h.getIntro)
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{
		service: newService(db, rdb),
	}
}

func (h *Handler) changeConfig(ctx *gin.Context) {
	var cfg SitedataConfig
	if ctx.ShouldBindBodyWithJSON(&cfg) != nil {
		h.Fail(ctx, httpserver.ErrRequest)
		return
	}

	setConfig(cfg)
	logger.Info(
		"change config successfully",
		zap.String("model", "article"),
	)
}

func (h *Handler) getIntro(ctx *gin.Context) {
	data, err := h.service.getIntro(ctx)
	if err != nil {
		h.Fail(ctx, httpserver.ErrDBOp)
		return
	}
	h.Success(ctx, data)
}
