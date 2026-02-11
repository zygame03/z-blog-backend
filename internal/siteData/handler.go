package sitedata

import (
	"my_web/backend/internal/httpserver"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	httpserver.BaseHandler
	service *service
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api")
	r.GET("/data/intro", h.getIntro)
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{
		service: newService(db, rdb),
	}
}

func (h *Handler) getIntro(ctx *gin.Context) {
	data, err := h.service.getIntro(ctx)
	if err != nil {
		h.Fail(ctx, httpserver.ErrDBOp)
		return
	}
	h.Success(ctx, data)
}
