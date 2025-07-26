package main

import (
	"log"
	"net/http"
	"time"

	"guandan-world/backend/auth"
	"guandan-world/backend/game"
	"guandan-world/backend/handlers"
	"guandan-world/backend/room"
	"guandan-world/backend/websocket"

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

	// 初始化房间服务
	roomService := room.NewRoomService(authService)

	// 初始化 WebSocket 管理器
	wsManager := websocket.NewWSManager(authService, roomService)

	// 初始化房间处理器
	roomHandler := handlers.NewRoomHandler(roomService, authService)

	// 初始化游戏服务（保留以备将来使用）
	_ = game.NewGameService(wsManager)

	// 初始化游戏驱动服务
	driverService := game.NewDriverService(wsManager)
	gameDriverHandler := handlers.NewGameDriverHandler(driverService)

	// 启动 WebSocket 管理器
	go wsManager.Run()

	// 注册认证路由
	authHandler.RegisterRoutes(r)

	// API 路由组
	api := r.Group("/api")
	{
		// 需要认证的路由
		protected := api.Group("/", authHandler.JWTMiddleware())
		{
			// 房间管理路由
			roomRoutes := protected.Group("/rooms")
			{
				roomRoutes.POST("/create", roomHandler.CreateRoom)
				roomRoutes.POST("/join", roomHandler.JoinRoom)
				roomRoutes.POST("/leave", roomHandler.LeaveRoom)
				roomRoutes.GET("/", roomHandler.GetRooms)
				roomRoutes.GET("/:id", roomHandler.GetRoom)
				roomRoutes.POST("/:id/start", roomHandler.StartGame)
			}

			// 游戏驱动路由
			driverRoutes := protected.Group("/game/driver")
			{
				driverRoutes.POST("/start", gameDriverHandler.StartGameWithDriver)
				driverRoutes.POST("/play-decision", gameDriverHandler.SubmitPlayDecision)
				driverRoutes.POST("/tribute-select", gameDriverHandler.SubmitTributeSelection)
				driverRoutes.POST("/tribute-return", gameDriverHandler.SubmitReturnTribute)
				driverRoutes.GET("/status/:room_id", gameDriverHandler.GetGameStatus)
				driverRoutes.POST("/stop/:room_id", gameDriverHandler.StopGame)
			}
		}
	}

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		// Get token from query parameter
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}
		
		// Validate token and get user
		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		
		// Pass user ID to WebSocket handler
		if err := wsManager.HandleWebSocket(c.Writer, c.Request, user.ID); err != nil {
			log.Printf("WebSocket error: %v", err)
		}
	})

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
