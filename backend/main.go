package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"guandan-world/backend/auth"
	"guandan-world/backend/config"
	"guandan-world/backend/database"
	"guandan-world/backend/game"
	"guandan-world/backend/handlers"
	"guandan-world/backend/room"
	"guandan-world/backend/websocket"
	"guandan-world/pkg/log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := log.Init("./logs", log.LevelInfo); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init logger:", err)
		os.Exit(1)
	}
	defer log.Close()

	cfg := config.Load()

	if cfg.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	log.Info("database connected successfully")

	if err := db.RunMigrations(); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("database migrations completed")

	userRepo := auth.NewPostgresUserRepository(db.DB)
	tokenRepo := auth.NewPostgresTokenRepository(db.DB)

	authService := auth.NewAuthService(userRepo, tokenRepo, auth.AuthServiceConfig{
		AccessSecret:       cfg.JWT.AccessSecret,
		RefreshSecret:      cfg.JWT.RefreshSecret,
		AccessTokenExpiry:  cfg.JWT.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
	})
	authHandler := handlers.NewAuthHandler(authService)

	roomService := room.NewRoomService(authService)
	wsManager := websocket.NewWSManager(authService, roomService)
	driverService := game.NewDriverService(wsManager)
	gameDriverHandler := handlers.NewGameDriverHandler(driverService)

	wsManager.SetReconnectHandler(driverService)
	roomHandler := handlers.NewRoomHandler(roomService, authService, driverService, wsManager)

	go wsManager.Run()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler.RegisterRoutes(r)

	api := r.Group("/api")
	{
		rooms := api.Group("/rooms")
		rooms.Use(authHandler.JWTMiddleware())
		{
			rooms.GET("", roomHandler.GetRooms)
			rooms.POST("/create", roomHandler.CreateRoom)
			rooms.GET("/my", roomHandler.GetMyRoom)
			rooms.GET("/:id", roomHandler.GetRoom)
			rooms.POST("/:id/join", roomHandler.JoinRoom)
			rooms.POST("/:id/leave", roomHandler.LeaveRoom)
			rooms.POST("/:id/start", roomHandler.StartGame)
		}

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

	r.GET("/ws", func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}

		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if err := wsManager.HandleWebSocket(c.Writer, c.Request, user.ID); err != nil {
			log.Error("WebSocket error", "error", err)
		}
	})

	r.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := db.HealthCheck(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Info("server starting", "port", cfg.Server.Port, "tls", cfg.Server.TLSEnabled)
		var err error
		if cfg.Server.TLSEnabled && cfg.Server.TLSCert != "" && cfg.Server.TLSKey != "" {
			err = srv.ListenAndServeTLS(cfg.Server.TLSCert, cfg.Server.TLSKey)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server exited")
}
