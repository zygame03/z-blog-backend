package http

import (
	"my_web/backend/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BaseHandler 提供通用的响应方法，可以被具体 Handler 嵌入
type BaseHandler struct{}

func (h *BaseHandler) Success(c *gin.Context, data any) {
	ReturnSuccess(c, data)
}

func (h *BaseHandler) Fail(c *gin.Context, err error) {
	logger.Error(
		"failed",
		zap.Error(err),
	)
	ReturnResponse(c, mapError(err), "")
}

func (h *BaseHandler) Response(c *gin.Context, r Result, data any) {
	ReturnResponse(c, r, data)
}

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

func ReturnResponse(c *gin.Context, r Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, SuccessResult, data)
}

func ReturnFail(c *gin.Context, r Result) {
	ReturnResponse(c, r, "")
}

// error到result的映射
func mapError(err error) Result {
	switch {
	default:
		return ErrDefault
	}
}
