package user

import (
	"my_web/backend/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	response.BaseHandler
	service *service
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{
		service: newService(db, rdb),
	}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/user")
	r.GET("/profile/:id", h.getProfile)
}

func (h *Handler) getProfile(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.Fail(ctx, response.ErrRequest)
		return
	}

	data, err := h.service.getProfile(ctx, id)
	if err != nil {
		h.Fail(ctx, response.ErrDBOp)
		return
	}

	h.Success(ctx, data)
}
