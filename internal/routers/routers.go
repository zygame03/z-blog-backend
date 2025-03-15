package routers

import (
	"my_web/backend/internal/handler"
	"my_web/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

var (
	articleAPI      handler.Article
	adminArticleAPI handler.AdminArticle
	userAPI         handler.User
)

func RegisterHandlers(r *gin.Engine) {
	registerArticleHandler(r)
	registerUserHandler(r)
	registerAdminHandler(r)
}

func registerArticleHandler(r *gin.Engine) {
	article := r.Group("api/article")
	{
		article.GET("", articleAPI.GetArticles)
		article.GET("/hotArticles", articleAPI.GetHotArticles)
		article.GET("/:id", articleAPI.GetArticleDetail)
	}
}

func registerUserHandler(r *gin.Engine) {
	use := r.Group("api/user")
	{
		// use.GET("/:id", userAPI.GetUser)
		use.POST("/login", userAPI.Login)
		use.POST("/register", userAPI.Register)
		use.GET("/profile", userAPI.GetProfile)
	}
}

func registerAdminHandler(r *gin.Engine) {
	admin := r.Group("/api/admin")
	{
		admin.Use(middleware.JWTAuth())
		admin.GET("/article", adminArticleAPI.AdminGetArticle)
		admin.GET("/article/:id", adminArticleAPI.AdminGetArticle)
		admin.POST("/article", adminArticleAPI.AdminSaveOrUpdateArticle)
		admin.DELETE("/article/:id", adminArticleAPI.AdminDeleteArticle)
	}
}
