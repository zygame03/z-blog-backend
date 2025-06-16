package main

import (
	"context"
	"log"
	"my_web/backend/internal/article"
	"my_web/backend/internal/config"
	"my_web/backend/internal/home"
	"my_web/backend/internal/infra"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/middleware"
	"my_web/backend/internal/site"

	"my_web/backend/internal/stats"
	"my_web/backend/internal/user"

	_ "my_web/backend/docs"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title z-blog-backend API
// @version 1.0
// @description Blog backend API
// @host localhost:8080
// @BasePath /api
func main() {
	logger.InitLogger()
	// load config
	cfg, err := config.LoadStConfig("./config/", "config")
	if err != nil {
		logger.Fatal(
			"加载静态配置失败",
			zap.Error(err),
		)
	}
	if cfg.JwtKey == "" {
		logger.Fatal(
			"JWT key 为空",
		)
	}

	err = config.LoadDyConfig("./config/", "dyconfig")
	if err != nil {
		log.Fatalf("加载动态配置失败: %v", err)
	}

	// initialize
	middleware.SetJWTKey(cfg.JwtKey)

	db, err := infra.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	logger.Info("数据库初始化成功")

	rdb, err := infra.InitCache(&cfg.Redis)
	if err != nil {
		log.Fatalf("Redis 初始化失败: %v", err)
	}
	logger.Info("Redis 初始化成功")

	cron := infra.NewCron()

	articleService := article.NewService(db, rdb, config.GetArticleConfig)
	articleService.RegisterCron(cron)
	articleHandler := article.NewHandler(articleService)

	dataService := site.NewService(db, rdb, config.GetSiteDataConfig)
	dataHandler := site.NewHandler(dataService)

	userService := user.NewService(db, rdb, config.GetUserConfig)
	userHandler := user.NewHandler(userService)

	homeService := home.NewService(db, rdb, func() {})
	homeHandler := home.NewHandler(homeService)

	statsService := stats.NewService(db, rdb, config.GetStatsConfig)
	statsService.RegisterCron(cron)
	statsHandler := stats.NewHandler(statsService)

	routers := []infra.Router{
		articleHandler,
		dataHandler,
		userHandler,
		homeHandler,
		statsHandler,
	}

	// dependency injection
	srv := infra.NewHttpserver(
		&cfg.Httpserver,
		routers,
		middleware.ViewsCounter(statsService),
	)

	cron.Start()
	defer cron.Stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务停止", err)
	}

	log.Println("退出")
}
