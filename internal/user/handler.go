package user

import (
	"my_web/backend/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	response.BaseHandler
	serv *service
}

func NewHandler(serv *service) *Handler {
	return &Handler{
		serv: serv,
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

	data, err := h.serv.getProfile(ctx, id)
	if err != nil {
		h.Fail(ctx, response.ErrDBOp)
		return
	}

	h.Success(ctx, data)
}
