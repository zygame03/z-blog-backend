package home

import (
	"my_web/backend/internal/config"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	response.BaseHandler
	serv *Service
}

func NewHandler(serv *Service) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/home")
	r.GET("/config/all", h.getSchemaAll)
	r.GET("/config/module/:name", h.getSchemaWithModule)
	// r.POST("/config/all", h.loadConfigAll)
	r.POST("/config/module/:name", h.loadConfigWithModule)
}

func (h *Handler) getSchemaWithModule(ctx *gin.Context) {
	module := ctx.Param("name")
	h.Success(ctx, response.GetSchemaByModule(module))
}

func (h *Handler) getSchemaAll(ctx *gin.Context) {
	h.Success(ctx, response.GetSchemaAll())
}

func (h *Handler) loadConfigWithModule(ctx *gin.Context) {
	module := ctx.Param("name")
	conf := config.NewConfigByModule(module)
	err := ctx.ShouldBindBodyWithJSON(conf)
	if err != nil {
		logger.Error(
			"invalid type",
			zap.String("module", module),
		)
		h.Fail(ctx, response.ErrConfigType)
		return
	}

	err = config.SetConfigWithModule(module, conf)
	if err != nil {
		logger.Error(
			"invalid type",
			zap.String("module", module),
		)
		h.Fail(ctx, response.ErrConfigType)
		return
	}
	h.Success(ctx, "")
}
