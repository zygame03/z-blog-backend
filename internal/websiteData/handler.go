package websiteData

import (
	"my_web/backend/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	response.BaseHandler
	service *service
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api")
	r.GET("/data/intro", h.getIntro)
}

func NewHandler(db *gorm.DB, rdb *redis.Client, cfg func() *WebsiteDataConfig) *Handler {
	return &Handler{
		service: newService(db, rdb, cfg),
	}
}

func (h *Handler) getIntro(ctx *gin.Context) {
	data, err := h.service.getIntro(ctx)
	if err != nil {
		h.Fail(ctx, response.ErrDBOp)
		return
	}
	h.Success(ctx, data)
}
