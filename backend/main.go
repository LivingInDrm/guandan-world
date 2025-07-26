package main

import (
	"log"
	"net/http"
	"time"

	"guandan-world/backend/auth"
	"guandan-world/backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 初始化认证服务
	authService := auth.NewAuthService("your-secret-key-change-in-production", 24*time.Hour)
	authHandler := handlers.NewAuthHandler(authService)

	// 注册认证路由
	authHandler.RegisterRoutes(r)

	// 健康检查接口
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "pong",
		})
	})

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
