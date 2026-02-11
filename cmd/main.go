package main

import (
	"context"
	"log"
	"my_web/backend/internal/article"
	"my_web/backend/internal/config"
	"my_web/backend/internal/httpserver"
	"my_web/backend/internal/infra"
	"my_web/backend/internal/logger"
	sitedata "my_web/backend/internal/siteData"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.InitLogger()
	// Load config
	config, err := config.ReadConfig("config/", "config", "json")
	if err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	// Initialization
	db, err := infra.InitDatabase(&config.Database)

	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	rdb, err := infra.InitRedis(&config.Redis)
	if err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}

	ctx := context.Background()
	articleServ := article.NewArticleService(ctx, db, rdb)
	articleHandler := article.NewHandler(articleServ)

	sitedataHandler := sitedata.NewHandler(db, rdb)

	// Start the httpserver in goroutine
	srv := httpserver.NewHttpserver(
		&config.Httpserver,
		articleHandler,
		sitedataHandler,
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("服务已退出")
}
