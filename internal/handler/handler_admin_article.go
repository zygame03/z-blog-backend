package handler

import (
	"my_web/backend/internal/global"
	"my_web/backend/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminArticle struct{}

type ArticleReq struct {
	ID         int       `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Title      string    `json:"title"`               // 标题
	Desc       string    `json:"desc" gorm:"text"`    // 描述
	Content    string    `json:"content" gorm:"text"` // 正文
	AuthorName string    `json:"authorName"`          // 作者
	Views      uint      `json:"views"`               // 浏览数
	Tags       string    `json:"tags"`                // 标签（逗号分隔形式）
	Cover      string    `json:"cover"`               // 封面
	Status     uint      `json:"status"`              // 状态
}

func (*AdminArticle) AdminGetArticle(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}
	pagesize, err := strconv.Atoi(c.DefaultQuery("pagesize", "10"))
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	s := GetAdminArticleService()
	data, total, err := s.AdminGetArticlesByPage(page, pagesize)
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, PageResult[models.Article]{
		Page:  page,
		Size:  pagesize,
		Data:  data,
		Total: int(total),
	})
}

func (*AdminArticle) AdminSaveOrUpdateArticle(c *gin.Context) {
	var data ArticleReq
	err := c.ShouldBindJSON(&data)
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	article := models.Article{
		Title:      data.Title,
		Desc:       data.Desc,
		Content:    data.Content,
		AuthorName: data.AuthorName,
		Views:      data.Views,
		Tags:       data.Tags,
		Cover:      data.Cover,
		Status:     data.Status,
	}

	if data.ID > 0 {
		article.ID = data.ID
	}

	s := GetAdminArticleService()

	if data.ID > 0 {
		err = s.AdminSaveOrUpdateArticle(&article)
	} else {
		err = s.AdminSaveOrUpdateArticle(&article)
	}

	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, article)
}

func (*AdminArticle) AdminDeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	s := GetAdminArticleService()
	err := s.AdminDeleteArticle(id)
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// id=?&status=?
func (*AdminArticle) AdminChangeStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	status, err := strconv.Atoi(c.DefaultQuery("status", "0"))
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	s := GetAdminArticleService()

	err = s.AdminChangeStatus(id, status)
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	ReturnSuccess(c, "success")
}
