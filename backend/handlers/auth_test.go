package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"guandan-world/backend/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() (*gin.Engine, *AuthHandler, auth.AuthService) {
	gin.SetMode(gin.TestMode)
	
	authService := auth.NewAuthService("test-secret", time.Hour)
	authHandler := NewAuthHandler(authService)
	
	router := gin.New()
	authHandler.RegisterRoutes(router)
	
	return router, authHandler, authService
}

func TestAuthHandler_Register(t *testing.T) {
	router, _, _ := setupTestRouter()

	t.Run("successful registration", func(t *testing.T) {
		reqBody := RegisterRequest{
			Username: "testuser",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "testuser", response.User.Username)
		assert.NotEmpty(t, response.User.ID)
		assert.True(t, response.User.Online)
		assert.NotEmpty(t, response.Token.Token)
		assert.Equal(t, response.User.ID, response.Token.UserID)
	})

	t.Run("invalid request format", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("missing username", func(t *testing.T) {
		reqBody := RegisterRequest{
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("duplicate username", func(t *testing.T) {
		// First registration
		reqBody := RegisterRequest{
			Username: "duplicate",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Second registration with same username
		req, _ = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "username_exists", response.Error)
	})

	t.Run("password too short", func(t *testing.T) {
		reqBody := RegisterRequest{
			Username: "shortpass",
			Password: "12345",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "registration_failed", response.Error)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	router, _, authService := setupTestRouter()

	// Register a test user first
	_, err := authService.Register("logintest", "password123")
	require.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		reqBody := LoginRequest{
			Username: "logintest",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "logintest", response.User.Username)
		assert.True(t, response.User.Online)
		assert.NotEmpty(t, response.Token.Token)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		reqBody := LoginRequest{
			Username: "logintest",
			Password: "wrongpassword",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "authentication_failed", response.Error)
	})

	t.Run("non-existent user", func(t *testing.T) {
		reqBody := LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request format", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	router, _, authService := setupTestRouter()

	// Register and login a test user
	_, err := authService.Register("logouttest", "password123")
	require.NoError(t, err)

	token, err := authService.Login("logouttest", "password123")
	require.NoError(t, err)

	t.Run("successful logout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token.Token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Successfully logged out", response["message"])

		// Verify token is no longer valid
		_, err = authService.ValidateToken(token.Token)
		assert.Error(t, err)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "missing_token", response.Error)
	})

	t.Run("invalid token format", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "invalid_token_format", response.Error)
	})

	t.Run("invalid token", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Me(t *testing.T) {
	router, _, authService := setupTestRouter()

	// Register and login a test user
	user, err := authService.Register("metest", "password123")
	require.NoError(t, err)

	token, err := authService.Login("metest", "password123")
	require.NoError(t, err)

	t.Run("successful me request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+token.Token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]*auth.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		returnedUser := response["user"]
		assert.Equal(t, user.ID, returnedUser.ID)
		assert.Equal(t, user.Username, returnedUser.Username)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/auth/me", nil)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/auth/me", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthHandler_JWTMiddleware(t *testing.T) {
	router, authHandler, authService := setupTestRouter()

	// Add a protected route for testing
	protected := router.Group("/api/protected")
	protected.Use(authHandler.JWTMiddleware())
	protected.GET("/test", func(c *gin.Context) {
		user, _ := c.Get("user")
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"user":    user,
			"user_id": userID,
		})
	})

	// Register and login a test user
	_, err := authService.Register("middlewaretest", "password123")
	require.NoError(t, err)

	token, err := authService.Login("middlewaretest", "password123")
	require.NoError(t, err)

	t.Run("valid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/protected/test", nil)
		req.Header.Set("Authorization", "Bearer "+token.Token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["user"])
		assert.NotNil(t, response["user_id"])
	})

	t.Run("missing token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/protected/test", nil)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/protected/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("expired token", func(t *testing.T) {
		// Create a new router with short expiry service
		gin.SetMode(gin.TestMode)
		shortExpiryService := auth.NewAuthService("test-secret", time.Millisecond*100)
		shortExpiryHandler := NewAuthHandler(shortExpiryService)
		
		shortExpiryRouter := gin.New()
		shortExpiryHandler.RegisterRoutes(shortExpiryRouter)
		
		// Add protected route
		protected := shortExpiryRouter.Group("/api/protected")
		protected.Use(shortExpiryHandler.JWTMiddleware())
		protected.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		_, err := shortExpiryService.Register("expiredtest", "password123")
		require.NoError(t, err)

		expiredToken, err := shortExpiryService.Login("expiredtest", "password123")
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(time.Millisecond * 200)

		req, _ := http.NewRequest("GET", "/api/protected/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken.Token)
		
		w := httptest.NewRecorder()
		shortExpiryRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response ErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "token_expired", response.Error)
	})
}

func TestAuthHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	authService := auth.NewAuthService("test-secret", time.Hour)
	authHandler := NewAuthHandler(authService)
	
	router := gin.New()
	authHandler.RegisterRoutes(router)

	// Test that all routes are registered
	routes := router.Routes()
	
	expectedRoutes := []string{
		"POST /api/auth/register",
		"POST /api/auth/login", 
		"POST /api/auth/logout",
		"GET /api/auth/me",
	}

	for _, expectedRoute := range expectedRoutes {
		found := false
		for _, route := range routes {
			if route.Method+" "+route.Path == expectedRoute {
				found = true
				break
			}
		}
		assert.True(t, found, "Route %s should be registered", expectedRoute)
	}
}