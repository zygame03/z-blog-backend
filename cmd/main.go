package main

import (
	"context"
	"log"
	"my_web/backend/internal/article"
	"my_web/backend/internal/config"
	"my_web/backend/internal/infra"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/middleware"
	"my_web/backend/internal/stats"
	"my_web/backend/internal/user"
	"my_web/backend/internal/websiteData"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.InitLogger()
	// load config
	cfg, err := config.LoadStConfig("./config/", "config")
	if err != nil {
		log.Fatalf("load static config failed: %v", err)
	}

	err = config.LoadDyConfig("./config/", "dyconfig")
	if err != nil {
		log.Fatalf("load dynamic config failed: %v", err)
	}

	// initialize
	middleware.SetJWTKey(cfg.JwtKey)

	db, err := infra.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("database initialized failed: %v", err)
	}

	rdb, err := infra.InitCache(&cfg.Redis)
	if err != nil {
		log.Fatalf("redis initialized failed: %v", err)
	}

	cron := infra.NewCron()
	ctx := context.Background()

	articleService := article.NewService(
		db, rdb,
		func() *article.ArticleConfig {
			return &config.GetConfig().Article
		},
	)
	articleService.RegisterCron(cron)
	articleHandler := article.NewHandler(articleService)

	dataHandler := websiteData.NewHandler(
		db, rdb,
		func() *websiteData.WebsiteDataConfig {
			return &config.GetConfig().SiteData
		},
	)
	userHandler := user.NewHandler(
		db, rdb,
	)

	routers := []infra.Router{
		articleHandler,
		dataHandler,
		userHandler,
	}

	statsService := stats.NewService(db, rdb, cron, ctx)

	// reject
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
	log.Println("close...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("exit")
}
