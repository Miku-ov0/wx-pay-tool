package main

import (
	"log"
	"os"
	"path/filepath"

	"wx-mch-trans/internal/handler"
	"wx-mch-trans/internal/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	// 创建必要的目录
	dirs := []string{"data/trans", "logs"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
	}

	// 初始化日志
	logFile, err := os.OpenFile(
		filepath.Join("logs", "app.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("创建日志文件失败: %v", err)
	}

	utils.InitLogger(logFile)
}

func main() {
	r := gin.Default()

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 设置静态文件路由
	r.Static("/static", "./static")

	// 注册路由
	handler.RegisterRoutes(r)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
