package admin

import (
	"my_web/backend/internal/middleware"
	"my_web/backend/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	response.BaseHandler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/admin", middleware.JWTAuth())
	r.GET("/config/remote/all", h.RemoteSendConfigAll)
	r.GET("/config/remote/module/:name", h.RemoteSendConfigWithModule)
	r.POST("/config/remote/all", h.RemoteLoadConfigAll)
	r.POST("/config/remote/module/:name", h.RemoteLoadConfigWithModule)
}

func (h *Handler) RemoteSendConfigWithModule(ctx *gin.Context) {

}

func (h *Handler) RemoteSendConfigAll(ctx *gin.Context) {

}

func (h *Handler) RemoteLoadConfigWithModule(ctx *gin.Context) {

}

func (h *Handler) RemoteLoadConfigAll(ctx *gin.Context) {
}
