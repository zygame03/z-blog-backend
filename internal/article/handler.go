package article

import (
	"my_web/backend/internal/http"
	"my_web/backend/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	http.BaseHandler
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
		r.POST("/comment", h.sendComment)

		r.Use(middleware.JWTAuth())
		r.POST("", h.saveArticle)
		r.POST("/comment/review", h.reviewComment)
	}
}

type ArticleListByPageResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    ArticlePageVO `json:"data"`
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
		h.Fail(ctx, err)
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	data, err := h.service.getArticlesByPage(ctx.Request.Context(), page, pageSize)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, data)
}

type ArticleListResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    []ArticleSummary `json:"data"`
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
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, data)
}

type ArticleDetailResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    ArticleDetailVO `json:"data"`
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
		h.Fail(ctx, err)
		return
	}

	// 获取userID
	userID := ctx.GetString("userID")
	if userID == "" {
		userID = ctx.ClientIP()
	}

	// 获取文章和首屏评论
	data, err := h.service.getArticleDetailWithComments(ctx.Request.Context(), id, userID, 1, 100)
	if err != nil {
		h.Fail(ctx, err)
		return
	}
	if data == nil {
		h.Response(ctx, http.ArticleNotFound, "")
		return
	}

	h.Success(ctx, data)
}

func (h *Handler) saveArticle(ctx *gin.Context) {
	var article Article

	err := ctx.ShouldBindJSON(&article)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	id, err := h.service.save(ctx, &article)
	if err != nil {
		h.Fail(ctx, err)
		return
	}
	h.Success(ctx, id)
}

func (h *Handler) sendComment(ctx *gin.Context) {
	var req struct {
		ArticleID int    `json:"article_id"`
		Username  string `json:"username"`
		Content   string `json:"content"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	id, err := h.service.sendComment(ctx.Request.Context(), SendCommentReq{
		ArticleID: req.ArticleID,
		Username:  req.Username,
		UserIP:    ctx.ClientIP(),
		Content:   req.Content,
	})
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, id)
}

func (h *Handler) reviewComment(ctx *gin.Context) {
	var req struct {
		ID     int `json:"id"`
		Status int `json:"Status"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	err = h.service.reviewComment(ctx, req.ID, CommentStatus(req.Status))
	if err != nil {
		h.Fail(ctx, err)
		return
	}

	h.Success(ctx, "")
}
