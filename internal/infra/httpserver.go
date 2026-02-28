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
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	ExposeHeaders    []string `mapstructure:"exposeHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"`
}

type Router interface {
	RegisterRoutes(*gin.Engine)
}

func NewHttpserver(cfg *HttpserverConfig, routers ...Router) *http.Server {
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

	for _, r := range routers {
		r.RegisterRoutes(e)
	}

	return &http.Server{
		Addr:    cfg.Port,
		Handler: e,
	}
}
