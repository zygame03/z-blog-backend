package handler

import (
	"my_web/backend/internal/constants"
	"my_web/backend/internal/global"
	"my_web/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type PageResult[T any] struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
	Data  []T `json:"data"`
}

func ReturnHttpResponse(c *gin.Context, httpcode, code int, msg string, data any) {
	c.JSON(httpcode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func ReturnResponse(c *gin.Context, r global.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, global.SuccessResult, data)
}

func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(constants.DATABASE).(*gorm.DB)
}

func GetRedis(c *gin.Context) *redis.Client {
	return c.MustGet(constants.REDIS).(*redis.Client)
}

// GetArticleService 获取文章服务实例（单例）
func GetArticleService() *service.ArticleService {
	return global.Services.ArticleService
}

// GetUserService 获取用户服务实例（单例）
func GetUserService() *service.UserService {
	return global.Services.UserService
}

// GetAdminArticleService 获取管理端文章服务实例（单例）
func GetAdminArticleService() *service.AdminArticleService {
	return global.Services.AdminArticleService
}
