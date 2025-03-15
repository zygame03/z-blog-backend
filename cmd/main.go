package main

import (
	"log"
	"my_web/backend/internal/constants"
	"my_web/backend/internal/global"
	"my_web/backend/internal/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 读取配置
	config, err := global.ReadConfig(constants.CONFIG_PATH)
	if err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	// 初始化应用依赖
	if err := global.Init(config, r); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 注册路由
	routers.RegisterHandlers(r)

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
