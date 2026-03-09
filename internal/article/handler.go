package article

import (
	"my_web/backend/internal/logger"
	"my_web/backend/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	response.BaseHandler
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/article")
	{
		r.GET("", h.getArticles)
		r.GET("/hotArticles", h.getHotArticles)
		r.GET("/:id", h.getArticleDetail)
	}
}

// 获取文章列表
func (h *Handler) getArticles(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		logger.Error(
			"parse query failed",
			zap.Error(err),
		)
		h.Fail(ctx, response.ErrRequest)
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil {
		logger.Error(
			"parse query failed",
			zap.Error(err),
		)
		h.Fail(ctx, response.ErrRequest)
		return
	}

	articles, total, err := h.service.GetArticlesByPage(ctx.Request.Context(), page, pageSize)
	if err != nil {
		logger.Error(
			"get article failed",
			zap.Int("page", page),
			zap.Int("pagesize", pageSize),
			zap.Error(err),
		)
		h.Fail(ctx, response.ErrDBOp)
		return
	}

	h.Success(ctx, response.PageResult[ArticleWithoutContent]{
		Page:  page,
		Size:  pageSize,
		Total: total,
		Data:  articles,
	})
}

// get hot articles by views
func (h *Handler) getHotArticles(ctx *gin.Context) {
	data, err := h.service.GetArticlesByPopular(ctx, 10)
	if err != nil {
		logger.Error(
			"get article failed",
			zap.Error(err),
		)
		h.Fail(ctx, response.ErrDBOp)
		return
	}

	h.Success(ctx, data)
}

// get an article with content
func (h *Handler) getArticleDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Error(
			"request invalid id",
			zap.Error(err),
		)
		h.Fail(ctx, response.ErrRequest)
		return
	}

	// get userid
	userID := ctx.GetString("userID") // if middleware set userud
	if userID == "" {
		userID = ctx.ClientIP() // use ip as userid
	}

	data, err := h.service.GetArticleByID(ctx.Request.Context(), id, userID)
	if err != nil {
		logger.Error(
			"get article failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		h.Fail(ctx, response.ErrDBOp)
		return
	}

	h.Success(ctx, data)
}
