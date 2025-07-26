package handlers

import (
	"net/http"
	"strings"

	"guandan-world/backend/auth"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService auth.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents a successful authentication response
type AuthResponse struct {
	User  *auth.User      `json:"user"`
	Token *auth.AuthToken `json:"token"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	// Register user
	_, err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "registration_failed"
		
		// Handle specific error cases
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

	// Auto-login after successful registration
	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		// Registration succeeded but login failed - this shouldn't happen
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "auto_login_failed",
			Message: "Registration succeeded but auto-login failed",
		})
		return
	}

	// Get updated user (should be online after login)
	updatedUser, err := h.authService.GetUserByID(token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "user_lookup_failed",
			Message: "Failed to retrieve updated user details",
		})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		User:  updatedUser,
		Token: token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	// Authenticate user
	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "authentication_failed",
			Message: err.Error(),
		})
		return
	}

	// Get user details
	user, err := h.authService.GetUserByID(token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "user_lookup_failed",
			Message: "Failed to retrieve user details",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User:  user,
		Token: token,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_token",
			Message: "Authorization header is required",
		})
		return
	}

	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_token_format",
			Message: "Token must be in 'Bearer <token>' format",
		})
		return
	}

	// Logout user
	err := h.authService.Logout(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "logout_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// Me returns the current user's information
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user from context (set by JWT middleware)
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

// JWTMiddleware validates JWT tokens and sets user context
func (h *AuthHandler) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "missing_token",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_token_format",
				Message: "Token must be in 'Bearer <token>' format",
			})
			c.Abort()
			return
		}

		// Validate token
		user, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			statusCode := http.StatusUnauthorized
			errorCode := "invalid_token"
			
			if strings.Contains(err.Error(), "token expired") {
				errorCode = "token_expired"
			}
			
			c.JSON(statusCode, ErrorResponse{
				Error:   errorCode,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// RegisterRoutes registers all authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.GET("/me", h.JWTMiddleware(), h.Me)
	}
}