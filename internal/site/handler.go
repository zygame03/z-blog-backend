package site

import (
	"my_web/backend/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	response.BaseHandler
	service *Service
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api")
	r.GET("/data/intro", h.getIntro)
	r.GET("/data/announcement", h.getAnnouncement)
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
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

func (h *Handler) getAnnouncement(ctx *gin.Context) {
	data, err := h.service.getAnnouncement(ctx)
	if err != nil {
		h.Fail(ctx, response.ErrDBOp)
		return
	}
	h.Success(ctx, data)
}
