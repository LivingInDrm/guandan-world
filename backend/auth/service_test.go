package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthService(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)
	assert.NotNil(t, service)
}

func TestAuthService_Register(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)

	t.Run("successful registration", func(t *testing.T) {
		user, err := service.Register("testuser", "password123")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		assert.NotEmpty(t, user.ID)
		assert.False(t, user.Online)
		assert.Empty(t, user.Password) // Password should not be exposed in returned user
	})

	t.Run("empty username", func(t *testing.T) {
		_, err := service.Register("", "password123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username cannot be empty")
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := service.Register("testuser2", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be empty")
	})

	t.Run("password too short", func(t *testing.T) {
		_, err := service.Register("testuser3", "12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must be at least 6 characters")
	})

	t.Run("duplicate username", func(t *testing.T) {
		// First registration should succeed
		_, err := service.Register("duplicate", "password123")
		require.NoError(t, err)

		// Second registration with same username should fail
		_, err = service.Register("duplicate", "password456")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
	})
}

func TestAuthService_Login(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)

	// Register a test user first
	user, err := service.Register("logintest", "password123")
	require.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		token, err := service.Login("logintest", "password123")
		require.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.Token)
		assert.Equal(t, user.ID, token.UserID)
		assert.True(t, token.ExpiresAt.After(time.Now()))

		// Check that user is marked as online
		loggedInUser, err := service.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.True(t, loggedInUser.Online)
	})

	t.Run("invalid username", func(t *testing.T) {
		_, err := service.Login("nonexistent", "password123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid username or password")
	})

	t.Run("invalid password", func(t *testing.T) {
		_, err := service.Login("logintest", "wrongpassword")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid username or password")
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)

	// Register and login a test user
	user, err := service.Register("validatetest", "password123")
	require.NoError(t, err)

	token, err := service.Login("validatetest", "password123")
	require.NoError(t, err)

	t.Run("valid token", func(t *testing.T) {
		validatedUser, err := service.ValidateToken(token.Token)
		require.NoError(t, err)
		assert.NotNil(t, validatedUser)
		assert.Equal(t, user.ID, validatedUser.ID)
		assert.Equal(t, user.Username, validatedUser.Username)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := service.ValidateToken("invalid-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("malformed token", func(t *testing.T) {
		_, err := service.ValidateToken("malformed.jwt.token")
		assert.Error(t, err)
	})
}

func TestAuthService_ValidateToken_Expired(t *testing.T) {
	// Create service with very short token expiry
	service := NewAuthService("test-secret", time.Millisecond*100)

	// Register and login a test user
	_, err := service.Register("expiredtest", "password123")
	require.NoError(t, err)

	token, err := service.Login("expiredtest", "password123")
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 200)

	// Token should now be expired
	_, err = service.ValidateToken(token.Token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token expired")
}

func TestAuthService_Logout(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)

	// Register and login a test user
	user, err := service.Register("logouttest", "password123")
	require.NoError(t, err)

	token, err := service.Login("logouttest", "password123")
	require.NoError(t, err)

	// Verify user is online
	loggedInUser, err := service.GetUserByID(user.ID)
	require.NoError(t, err)
	assert.True(t, loggedInUser.Online)

	t.Run("successful logout", func(t *testing.T) {
		err := service.Logout(token.Token)
		require.NoError(t, err)

		// Token should no longer be valid
		_, err = service.ValidateToken(token.Token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")

		// User should be marked as offline
		loggedOutUser, err := service.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.False(t, loggedOutUser.Online)
	})

	t.Run("logout with invalid token", func(t *testing.T) {
		err := service.Logout("invalid-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})
}

func TestAuthService_GetUserByID(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)

	// Register a test user
	user, err := service.Register("getusertest", "password123")
	require.NoError(t, err)

	t.Run("existing user", func(t *testing.T) {
		foundUser, err := service.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := service.GetUserByID("non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestAuthService_ConcurrentAccess(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)

	// Test concurrent registration
	t.Run("concurrent registration", func(t *testing.T) {
		done := make(chan bool, 10)
		errors := make(chan error, 10)

		for i := 0; i < 10; i++ {
			go func(i int) {
				username := fmt.Sprintf("concurrent%d", i)
				_, err := service.Register(username, "password123")
				if err != nil {
					errors <- err
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Check if there were any errors
		select {
		case err := <-errors:
			t.Errorf("Concurrent registration failed: %v", err)
		default:
			// No errors, test passed
		}
	})
}

func TestAuthService_TokenSecurity(t *testing.T) {
	service1 := NewAuthService("secret1", time.Hour)
	service2 := NewAuthService("secret2", time.Hour)

	// Register user in service1
	_, err := service1.Register("securitytest", "password123")
	require.NoError(t, err)

	// Login with service1
	token, err := service1.Login("securitytest", "password123")
	require.NoError(t, err)

	// Token from service1 should not be valid in service2 (different secret)
	_, err = service2.ValidateToken(token.Token)
	assert.Error(t, err)
}