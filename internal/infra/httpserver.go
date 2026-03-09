package infra

import (
	"my_web/backend/internal/logger"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpserverConfig struct {
	Port string      `mapstructure:"port"`
	Cors *CorsConfig `mapstructure:"cors"`
}

// CorsConfig CORS 配置
type CorsConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

type Router interface {
	RegisterRoutes(*gin.Engine)
}

func NewHttpserver(cfg *HttpserverConfig, routers []Router, opts ...gin.HandlerFunc) *http.Server {
	e := gin.New()

	e.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Cors.AllowedOrigins,
		AllowMethods:     cfg.Cors.AllowedMethods,
		AllowHeaders:     cfg.Cors.AllowedHeaders,
		ExposeHeaders:    cfg.Cors.ExposeHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		MaxAge:           time.Duration(cfg.Cors.MaxAge) * time.Hour,
	}))
	e.Use(logger.GinLogger())

	for _, opt := range opts {
		e.Use(opt)
	}

	for _, r := range routers {
		r.RegisterRoutes(e)
	}

	return &http.Server{
		Addr:    cfg.Port,
		Handler: e,
	}
}
