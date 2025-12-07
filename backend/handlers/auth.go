package handlers

import (
	"net/http"
	"strings"

	"guandan-world/backend/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User  *auth.User       `json:"user"`
	Token *auth.TokenPair  `json:"token"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	_, err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "registration_failed"
		
		if strings.Contains(err.Error(), "username already exists") {
			statusCode = http.StatusConflict
			errorCode = "username_exists"
		}
		
		c.JSON(statusCode, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	tokenPair, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "auto_login_failed",
			Message: "Registration succeeded but auto-login failed",
		})
		return
	}

	updatedUser, err := h.authService.GetUserByID(tokenPair.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "user_lookup_failed",
			Message: "Failed to retrieve updated user details",
		})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		User:  updatedUser,
		Token: tokenPair,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	tokenPair, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "authentication_failed",
			Message: err.Error(),
		})
		return
	}

	user, err := h.authService.GetUserByID(tokenPair.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "user_lookup_failed",
			Message: "Failed to retrieve user details",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User:  user,
		Token: tokenPair,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "refresh_token is required",
		})
		return
	}

	tokenPair, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "refresh_failed",
			Message: err.Error(),
		})
		return
	}

	user, err := h.authService.GetUserByID(tokenPair.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "user_lookup_failed",
			Message: "Failed to retrieve user details",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User:  user,
		Token: tokenPair,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	if req.RefreshToken != "" {
		if err := h.authService.Logout(req.RefreshToken); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "logout_failed",
				Message: err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user data in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *AuthHandler) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "missing_token",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_token_format",
				Message: "Token must be in 'Bearer <token>' format",
			})
			c.Abort()
			return
		}

		user, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			errorCode := "invalid_token"
			
			if strings.Contains(err.Error(), "token is expired") {
				errorCode = "token_expired"
			}
			
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   errorCode,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.Refresh)
		authGroup.POST("/logout", h.Logout)
		authGroup.GET("/me", h.JWTMiddleware(), h.Me)
	}
}
