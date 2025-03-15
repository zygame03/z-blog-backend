package handler

import (
	"my_web/backend/internal/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct{}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (*User) Login(c *gin.Context) {
	var req LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	s := GetUserService()
	user, err := s.Login(req.Username, req.Password)
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":  user,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("your_secret_key"))

	ReturnSuccess(c, tokenString)
}

func (*User) Register(c *gin.Context) {
	var req LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnResponse(c, global.ErrRequest, err)
		return
	}

	s := GetUserService()
	user, err := s.Register(req.Username, req.Password)
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
	}

	ReturnSuccess(c, user)
}

func (*User) GetUser(c *gin.Context) {
	s := GetUserService()

	data, err := s.GetUser()
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, data)
}

func (*User) GetProfile(c *gin.Context) {
	s := GetUserService()

	data, err := s.GetProfile()
	if err != nil {
		ReturnResponse(c, global.ErrDBOp, err)
		return
	}

	ReturnSuccess(c, data)
}
