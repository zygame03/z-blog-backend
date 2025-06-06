package middleware

import (
	"my_web/backend/internal/stats"

	"github.com/gin-gonic/gin"
)

func ViewsCounter(s *stats.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		ctx.Next()

		// 结果校验
		if ctx.Writer.Status() == 200 {
			s.RecordUV(ip)
		}
	}
}
