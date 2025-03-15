package models

import "gorm.io/gorm"

const (
	Public = iota
	Private
)

type Article struct {
	Model
	Title      string `json:"title"`               // 标题
	Desc       string `json:"desc" gorm:"text"`    // 描述
	Content    string `json:"content" gorm:"text"` // 正文
	AuthorName string `json:"authorName"`          // 作者
	Views      uint   `json:"views"`               // 浏览数
	Tags       string `json:"tags"`                // 标签（逗号分隔形式）
	Cover      string `json:"cover"`               // 封面
	Status     uint   `json:"status"`              // 状态
	IsDelete   bool   `json:"is_delete"`
}

func (a *Article) Create(db *gorm.DB) error {
	return db.Create(a).Error
}

func (a *Article) Update(db *gorm.DB) error {
	return db.Model(a).Where("id = ?", a.ID).Updates(a).Error
}

func (a *Article) Delete(db *gorm.DB) error {
	return db.Model(a).Where("id = ?", a.ID).Update("is_delete", true).Error
}
