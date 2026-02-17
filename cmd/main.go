package main

import (
	"context"
	"log"
	"my_web/backend/internal/article"
	"my_web/backend/internal/config"
	"my_web/backend/internal/data"
	"my_web/backend/internal/httpserver"
	"my_web/backend/internal/infra"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/user"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.InitLogger()
	// 读取配置
	config, err := config.ReadConfig("config/", "config", "json")
	if err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	// 初始化应用依赖
	db, err := infra.InitDatabase(&config.Database)

	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	rdb, err := infra.InitRedis(&config.Redis)
	if err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}

	ctx := context.Background()
	articleHandler := article.NewHandler(ctx, db, rdb)
	sitedataHandler := data.NewHandler(db, rdb)
	userHandler := user.NewHandler(db, rdb)

	// 在 goroutine 中启动服务
	srv := httpserver.NewHttpserver(
		&config.Httpserver,
		articleHandler,
		sitedataHandler,
		userHandler,
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 优雅退出处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("close...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("exit")
}
