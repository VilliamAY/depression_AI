package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"depression_go/inits"
	routers "depression_go/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// 加载环境变量 - 指定config.env文件
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("未找到config.env文件，尝试加载.env文件")
		if err := godotenv.Load(); err != nil {
			log.Println("未找到.env文件，使用系统环境变量")
		}
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化数据库
	inits.InitDatabase()

	// 初始化JWT
	inits.InitConfig()

	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 设置路由
	routers.SetupRoutes(r)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动在端口 %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	// 关闭数据库连接
	inits.CloseDatabase()

	log.Println("服务器已关闭")
}
