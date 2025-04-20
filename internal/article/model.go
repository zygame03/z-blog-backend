package article

import (
	"my_web/backend/internal/global"
)

type ArticleStatus uint8

const (
	ArticlePublic = iota
	ArticlePrivate
)

type Article struct {
	global.BaseModel
	Title      string        `json:"title"`               // 标题
	Desc       string        `json:"desc" gorm:"text"`    // 描述
	Content    string        `json:"content" gorm:"text"` // 正文
	AuthorName string        `json:"authorName"`          // 作者
	Views      uint          `json:"views"`               // 浏览数
	Tags       string        `json:"tags"`                // 标签（逗号分隔形式）
	Cover      string        `json:"cover"`               // 封面
	Status     ArticleStatus `json:"status"`              // 状态
	IsDelete   bool          `json:"is_delete"`
}

// ---------------------------------------

type ArticleWithoutContent struct {
	global.BaseModel
	Title      string `json:"title"`
	AuthorName string `json:"authorName"` // 作者
	Views      uint   `json:"views"`      // 浏览数
	Tags       string `json:"tags"`       // 标签（逗号分隔）
	Cover      string `json:"cover"`      // 封面
}
