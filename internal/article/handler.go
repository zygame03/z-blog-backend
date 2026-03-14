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
	service *ArticleService
}

func NewHandler(service *ArticleService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(e *gin.Engine) {
	r := e.Group("/api/article")
	{
		r.GET("", h.getArticles)
		r.GET("/hot_articles", h.getHotArticles)
		r.GET("/:id", h.getArticleDetail)
	}
}

type ArticleListByPageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Page  int                     `json:"page"`
		Size  int                     `json:"size"`
		Total int                     `json:"total"`
		Data  []ArticleWithoutContent `json:"data"`
	} `json:"data"`
}

// 获取文章列表
// @Summary 获取文章列表
// @Description 分页获取文章
// @Tags article
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} ArticleListByPageResponse
// @Router /article [get]
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

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		logger.Error(
			"parse query failed",
			zap.Error(err),
		)
		h.Fail(ctx, response.ErrRequest)
		return
	}

	articles, total, err := h.service.getArticlesByPage(ctx.Request.Context(), page, pageSize)
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

type ArticleListResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    []ArticleWithoutContent `json:"data"`
}

// 获取热门文章
// @Summary 获取热门文章
// @Description 获取热门文章
// @Tags article
// @Accept json
// @Produce json
// @Success 200 {object} ArticleListResponse
// @Router /article/hot_articles [get]
func (h *Handler) getHotArticles(ctx *gin.Context) {
	data, err := h.service.getArticlesByPopular(ctx, 10)
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

type ArticleDetailResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    Article `json:"data"`
}

// 获取文章详情（带正文）
// @Summary 获取文章详情
// @Description 根据id获取文章详情（带正文）
// @Tags article
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} ArticleDetailResponse
// @Router /article/{id} [get]
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

	// 获取userID
	userID := ctx.GetString("userID")
	if userID == "" {
		userID = ctx.ClientIP()
	}

	data, err := h.service.getArticleByID(ctx.Request.Context(), id, userID)
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
