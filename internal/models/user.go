package models

import "gorm.io/gorm"

type User struct {
	Model
	Avatar      string `json:"avatar"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	Signature   string `json:"signature"`
	Github      string `json:"github"`
	Bilibili    string `json:"bilibili"`
	Skills      string `json:"skills"`
	Hobbies     string `json:"hobbies"`
	Timeline    string `json:"timeline"`
	FutureGoals string `json:"futureGoals"`
	Intro       string `json:"intro"`
}

func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) Update(db *gorm.DB) error {
	return db.Model(u).Where("id = ?", u.ID).Updates(u).Error
}

func (u *User) Delete(db *gorm.DB) error {
	return db.Model(u).Where("id = ?", u.ID).Update("is_delete", true).Error
}
