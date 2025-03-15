package handler

import (
	"my_web/backend/internal/global"
	"my_web/backend/internal/repo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Article struct {
}

// 获取文章列表
func (*Article) GetArticles(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	s := GetArticleService()
	articles, total, err := s.GetArticlesByPage(c.Request.Context(), page, pageSize)
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, PageResult[repo.ArticleVO]{
		Page:  page,
		Size:  pageSize,
		Total: total,
		Data:  articles,
	})
}

func (*Article) GetHotArticles(c *gin.Context) {
	s := GetArticleService()

	data, err := s.GetArticlesByPopular(10)
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, data)
}

// 获取文章详情（带正文）
func (*Article) GetArticleDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	s := GetArticleService()
	data, err := s.GetArticleByID(id)
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, data)
}
