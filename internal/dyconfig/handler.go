package dyconfig

import (
	"my_web/backend/internal/httpserver"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	httpserver.BaseHandler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/dyconfig")
	r.GET("/remote/config/all", h.RemoteSendConfigAll)
	r.GET("/remote/config/:name", h.RemoteSendConfigWithModule)
}

func (h *Handler) RemoteSendConfigWithModule(ctx *gin.Context) {

}

func (h *Handler) RemoteSendConfigAll(ctx *gin.Context) {

}

func (h *Handler) RemoteLoadConfigWithModule(ctx *gin.Context) {

}

func (h *Handler) RemoteLoadConfigAll(ctx *gin.Context) {
}
