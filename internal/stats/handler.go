package stats

import (
	"my_web/backend/internal/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	http.BaseHandler
	serv *Service
}

func NewHandler(serv *Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/stats")
	r.GET("/views", h.getViews)
}

func (h *Handler) getViews(ctx *gin.Context) {
	data, err := h.serv.getViews(ctx)
	if err != nil {
		h.Fail(ctx, err)
		return
	}
	h.Success(ctx, data)
}
