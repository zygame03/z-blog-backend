package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	log *zap.Logger
)

func InitLogger() {
	l, err := zap.NewDevelopment(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic("logger initialized failed")
	}

	log = l
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.FullPath()
		if path == "" {
			path = ctx.Request.URL.Path
		}
		method := ctx.Request.Method
		clientIP := ctx.ClientIP()

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()

		Info("http request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)

		// 异常日志打印
		// 自动映射？
	}
}
