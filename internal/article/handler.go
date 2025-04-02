package article

import (
	"context"
	"my_web/backend/internal/httpserver"
	"my_web/backend/internal/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	httpserver.BaseHandler
	service *service
}

func NewHandler(ctx context.Context, db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{
		service: newArticleService(ctx, db, rdb),
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
		h.Fail(ctx, httpserver.ErrRequest)
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil {
		logger.Error(
			"parse query failed",
			zap.Error(err),
		)
		h.Fail(ctx, httpserver.ErrRequest)
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
		h.Fail(ctx, httpserver.ErrDBOp)
		return
	}

	h.Success(ctx, httpserver.PageResult[ArticleWithoutContent]{
		Page:  page,
		Size:  pageSize,
		Total: total,
		Data:  articles,
	})
}

func (h *Handler) getHotArticles(ctx *gin.Context) {
	data, err := h.service.GetArticlesByPopular(ctx, 10)
	if err != nil {
		logger.Error(
			"get article failed",
			zap.Error(err),
		)
		h.Fail(ctx, httpserver.ErrDBOp)
		return
	}

	h.Success(ctx, data)
}

// 获取文章详情（带正文）
func (h *Handler) getArticleDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logger.Error(
			"request invalid id",
			zap.Error(err),
		)
		h.Fail(ctx, httpserver.ErrRequest)
		return
	}

	// 获取用户标识（优先使用用户ID，否则使用IP地址）
	userID := ctx.GetString("userID") // 如果中间件设置了用户ID
	if userID == "" {
		userID = ctx.ClientIP() // 使用IP地址作为标识
	}

	data, err := h.service.GetArticleByID(ctx.Request.Context(), id, userID)
	if err != nil {
		logger.Error(
			"get article failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		h.Fail(ctx, httpserver.ErrDBOp)
		return
	}

	h.Success(ctx, data)
}
