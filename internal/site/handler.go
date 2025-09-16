package site

import (
	"my_web/backend/internal/http"
	"my_web/backend/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	http.BaseHandler
	service *Service
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api")
	r.GET("/site/intro", h.getIntro)
	r.GET("/site/announcement", h.getAnnouncement)
	r.GET("/site/danmaku", h.getDanmaku)
	r.POST("/site/danmaku", h.sendDanmaku)

	admin := e.Group("/api/admin")
	admin.Use(middleware.JWTAuth())
	admin.GET("/site/danmaku", h.getDanmakuAdmin)
	admin.POST("/site/danmaku/:id/review", h.reviewDanmaku)
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) getIntro(ctx *gin.Context) {
	data, err := h.service.getIntro(ctx)
	if err != nil {
		h.Fail(ctx, err)
		return
	}
	h.Success(ctx, data)
}

func (h *Handler) getAnnouncement(ctx *gin.Context) {
	data, err := h.service.getAnnouncement(ctx)
	if err != nil {
		h.Fail(ctx, err)
		return
	}
	h.Success(ctx, data)
}

func (h *Handler) getDanmaku(ctx *gin.Context) {
	data, err := h.service.listDanmaku(ctx)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, data)
}

func (h *Handler) getDanmakuAdmin(ctx *gin.Context) {
	data, err := h.service.listDanmakuAdmin(ctx)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, data)
}

func (h *Handler) reviewDanmaku(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	var req struct {
		Status DanmakuStatus `json:"status"`
	}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	err = h.service.reviewDanmaku(ctx, id, req.Status)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, "")
}

func (h *Handler) sendDanmaku(ctx *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	senderID := ctx.GetString("userID")
	if senderID == "" {
		senderID = ctx.ClientIP()
	}
	ip := ctx.ClientIP()

	id, err := h.service.sendDanmaku(ctx.Request.Context(), req.Content, senderID, ip)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, gin.H{
		"id": id,
	})
}
