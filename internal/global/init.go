package global

import (
	"fmt"
	"log"
	"my_web/backend/internal/middleware"
	"my_web/backend/internal/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init 初始化应用依赖（数据库、Redis、CORS等）
func Init(cfg *Config, e *gin.Engine) error {
	if err := InitCors(cfg, e); err != nil {
		return fmt.Errorf("初始化 CORS 失败: %w", err)
	}

	// 初始化数据库
	db, err := InitDatabase(cfg, e)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 初始化 Redis
	rdb, err := InitRedis(cfg, e)
	if err != nil {
		return fmt.Errorf("初始化 Redis 失败: %w", err)
	}

	// 初始化服务容器
	InitServices(db, rdb)

	log.Println("所有组件初始化成功")
	return nil
}

// InitCors 初始化 CORS 中间件
func InitCors(cfg *Config, e *gin.Engine) error {
	conf := cfg.Cors
	e.Use(cors.New(cors.Config{
		AllowOrigins:     conf.AllowedOrigins,
		AllowMethods:     conf.AllowedMethods,
		AllowHeaders:     conf.AllowedHeaders,
		ExposeHeaders:    conf.ExposeHeaders,
		AllowCredentials: conf.AllowCredentials,
		MaxAge:           time.Duration(conf.MaxAge) * time.Hour,
	}))
	log.Println("CORS 初始化成功")
	return nil
}

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *Config, e *gin.Engine) (*gorm.DB, error) {
	conf := cfg.Database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port, conf.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&models.Article{},
		&models.User{},
	); err != nil {
		return nil, fmt.Errorf("数据库自动迁移失败: %w", err)
	}

	e.Use(middleware.WithGormDB(db))
	log.Println("数据库初始化成功")
	return db, nil
}

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *Config, e *gin.Engine) (*redis.Client, error) {
	conf := cfg.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
		Protocol: conf.Protocol,
	})

	e.Use(middleware.WithRedis(rdb))
	log.Println("Redis 初始化成功")
	return rdb, nil
}
