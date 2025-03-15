package global

import (
	"my_web/backend/internal/cache"
	"my_web/backend/internal/repo"
	"my_web/backend/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServiceContainer 服务容器，管理所有 service 实例（单例模式）
type ServiceContainer struct {
	// Services
	ArticleService      *service.ArticleService
	UserService         *service.UserService
	AdminArticleService *service.AdminArticleService

	// Repositories
	ArticleRepo *repo.ArticleRepo

	// Caches
	ArticleCache *cache.ArticleCache

	// Dependencies
	DB    *gorm.DB
	Redis *redis.Client
}

var (
	// Services 全局服务容器实例
	Services *ServiceContainer
)

// InitServices 初始化服务容器
func InitServices(db *gorm.DB, rdb *redis.Client) {
	// 初始化 Repositories
	articleRepo := &repo.ArticleRepo{DB: db}

	// 初始化 Caches
	articleCache := &cache.ArticleCache{RDB: rdb}

	// 初始化 Services
	articleService := service.NewArticleService(articleRepo, articleCache)
	userService := service.NewUserService(db)
	adminArticleService := service.NewAdminArticleService(db)

	Services = &ServiceContainer{
		ArticleService:      articleService,
		UserService:         userService,
		AdminArticleService: adminArticleService,
		ArticleRepo:         articleRepo,
		ArticleCache:        articleCache,
		DB:                  db,
		Redis:               rdb,
	}
}
