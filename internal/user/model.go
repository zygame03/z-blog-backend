package user

import (
	"errors"
	"my_web/backend/internal/global"
)

type User struct {
	global.BaseModel
	Username   string
	Password   string
	Permission int64
}

type Profile struct {
	global.BaseModel
	Name   string
	Desc   string
	Locate string
	Age    int
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserAlreadyExist = errors.New("user already exist")
)
