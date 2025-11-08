package main

import (
	"log"
	"net/http"
	"os"
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
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}
	authService := auth.NewAuthService(jwtSecret, 24*time.Hour)
	authHandler := handlers.NewAuthHandler(authService)

	// 初始化房间服务
	roomService := room.NewRoomService(authService)

	// 初始化 WebSocket 管理器
	wsManager := websocket.NewWSManager(authService, roomService)

	// 初始化游戏服务（保留以备将来使用）
	_ = game.NewGameService(wsManager)

	// 初始化游戏驱动服务
	driverService := game.NewDriverService(wsManager)
	gameDriverHandler := handlers.NewGameDriverHandler(driverService)

	// 初始化房间处理器 - 注入 driverService 和 wsManager
	roomHandler := handlers.NewRoomHandler(roomService, authService, driverService, wsManager)

	// 启动 WebSocket 管理器
	go wsManager.Run()

	// 注册认证路由（不需要JWT认证）
	authHandler.RegisterRoutes(r)

	// API 路由组 - 需要认证的路由
	api := r.Group("/api")
	{
		// 房间管理路由（需要认证）
		rooms := api.Group("/rooms")
		rooms.Use(authHandler.JWTMiddleware())
		{
			// 注意：特定路径必须在参数路径之前定义，避免被 /:id 匹配
			rooms.GET("", roomHandler.GetRooms)                // GET /api/rooms - list rooms
			rooms.POST("/create", roomHandler.CreateRoom)      // POST /api/rooms/create - create room  
			rooms.GET("/my", roomHandler.GetMyRoom)            // GET /api/rooms/my - get current user's room
			
			// 参数路径放在后面
			rooms.GET("/:id", roomHandler.GetRoom)             // GET /api/rooms/:id - get room details
			rooms.POST("/:id/join", roomHandler.JoinRoom)      // POST /api/rooms/:id/join - join room
			rooms.POST("/:id/leave", roomHandler.LeaveRoom)    // POST /api/rooms/:id/leave - leave room
			rooms.POST("/:id/start", roomHandler.StartGame)    // POST /api/rooms/:id/start - start game
		}

		// 游戏驱动路由（需要认证）
		driver := api.Group("/game/driver")
		driver.Use(authHandler.JWTMiddleware())
		{
			driver.POST("/start", gameDriverHandler.StartGameWithDriver)
			driver.POST("/play-decision", gameDriverHandler.SubmitPlayDecision)
			driver.POST("/tribute-select", gameDriverHandler.SubmitTributeSelection)
			driver.POST("/tribute-return", gameDriverHandler.SubmitReturnTribute)
			driver.GET("/status/:room_id", gameDriverHandler.GetGameStatus)
			driver.POST("/stop/:room_id", gameDriverHandler.StopGame)
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
