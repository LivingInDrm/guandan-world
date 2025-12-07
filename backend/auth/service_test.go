package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthService(t *testing.T) {
	service := NewTestAuthService("test-secret", time.Hour)
	assert.NotNil(t, service)
}

func TestAuthService_Register(t *testing.T) {
	service := NewTestAuthService("test-secret", time.Hour)

	t.Run("successful registration", func(t *testing.T) {
		user, err := service.Register("testuser", "Password123!")
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		assert.NotEmpty(t, user.ID)
		assert.False(t, user.Online)
		assert.Empty(t, user.Password) // Password should not be exposed in returned user
	})

	t.Run("username too short", func(t *testing.T) {
		_, err := service.Register("abc", "Password123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username must be at least 4 characters")
	})

	t.Run("username too long", func(t *testing.T) {
		_, err := service.Register("abcdefghijklmnopqrstu", "Password123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username must not exceed 20 characters")
	})

	t.Run("username with invalid characters", func(t *testing.T) {
		_, err := service.Register("user@name", "Password123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username can only contain letters, numbers, and underscores")
	})

	t.Run("username with chinese characters", func(t *testing.T) {
		_, err := service.Register("用户名test", "Password123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username can only contain letters, numbers, and underscores")
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := service.Register("testuser2", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be empty")
	})

	t.Run("password too short", func(t *testing.T) {
		_, err := service.Register("testuser3", "Pass1!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must be at least 8 characters")
	})

	t.Run("password too long", func(t *testing.T) {
		longPassword := "Password1!" + string(make([]byte, 50))
		_, err := service.Register("testuser3b", longPassword)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must not exceed 50 characters")
	})

	t.Run("password missing uppercase", func(t *testing.T) {
		_, err := service.Register("testuser4", "password123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least one uppercase letter")
	})

	t.Run("password missing lowercase", func(t *testing.T) {
		_, err := service.Register("testuser5", "PASSWORD123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least one lowercase letter")
	})

	t.Run("password missing digit", func(t *testing.T) {
		_, err := service.Register("testuser6", "Password!!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least one digit")
	})

	t.Run("password missing special character", func(t *testing.T) {
		_, err := service.Register("testuser7", "Password123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must contain at least one special character")
	})

	t.Run("duplicate username", func(t *testing.T) {
		// First registration should succeed
		_, err := service.Register("duplicate", "Password123!")
		require.NoError(t, err)

		// Second registration with same username should fail
		_, err = service.Register("duplicate", "Password456!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
	})
}

func TestAuthService_Login(t *testing.T) {
	service := NewTestAuthService("test-secret", time.Hour)

	// Register a test user first
	user, err := service.Register("logintest", "Password123!")
	require.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		token, err := service.Login("logintest", "Password123!")
		require.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token.AccessToken)
		assert.Equal(t, user.ID, token.UserID)
		assert.True(t, token.ExpiresAt.After(time.Now()))

		// Check that user is marked as online
		loggedInUser, err := service.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.True(t, loggedInUser.Online)
	})

	t.Run("invalid username", func(t *testing.T) {
		_, err := service.Login("nonexistent", "Password123!")
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
	service := NewTestAuthService("test-secret", time.Hour)

	// Register and login a test user
	user, err := service.Register("validatetest", "Password123!")
	require.NoError(t, err)

	token, err := service.Login("validatetest", "Password123!")
	require.NoError(t, err)

	t.Run("valid token", func(t *testing.T) {
		validatedUser, err := service.ValidateToken(token.AccessToken)
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
	service := NewTestAuthService("test-secret", time.Millisecond*100)

	// Register and login a test user
	_, err := service.Register("expiredtest", "Password123!")
	require.NoError(t, err)

	token, err := service.Login("expiredtest", "Password123!")
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 200)

	// Token should now be expired
	_, err = service.ValidateToken(token.AccessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestAuthService_Logout(t *testing.T) {
	service := NewTestAuthService("test-secret", time.Hour)

	// Register and login a test user
	user, err := service.Register("logouttest", "Password123!")
	require.NoError(t, err)

	token, err := service.Login("logouttest", "Password123!")
	require.NoError(t, err)

	// Verify user is online
	loggedInUser, err := service.GetUserByID(user.ID)
	require.NoError(t, err)
	assert.True(t, loggedInUser.Online)

	t.Run("successful logout", func(t *testing.T) {
		err := service.Logout(token.RefreshToken)
		require.NoError(t, err)

		// User should be marked as offline
		loggedOutUser, err := service.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.False(t, loggedOutUser.Online)
	})

	t.Run("logout with invalid token", func(t *testing.T) {
		err := service.Logout("invalid-token")
		assert.NoError(t, err) // Logout with invalid token returns nil
	})
}

func TestAuthService_GetUserByID(t *testing.T) {
	service := NewTestAuthService("test-secret", time.Hour)

	// Register a test user
	user, err := service.Register("getusertest", "Password123!")
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
	service := NewTestAuthService("test-secret", time.Hour)

	// Test concurrent registration
	t.Run("concurrent registration", func(t *testing.T) {
		done := make(chan bool, 10)
		errors := make(chan error, 10)

		for i := 0; i < 10; i++ {
			go func(i int) {
				username := fmt.Sprintf("concurrent%d", i)
				_, err := service.Register(username, "Password123!")
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
	service1 := NewTestAuthService("secret1", time.Hour)
	service2 := NewTestAuthService("secret2", time.Hour)

	// Register user in service1
	_, err := service1.Register("securitytest", "Password123!")
	require.NoError(t, err)

	// Login with service1
	token, err := service1.Login("securitytest", "Password123!")
	require.NoError(t, err)

	// Token from service1 should not be valid in service2 (different secret)
	_, err = service2.ValidateToken(token.AccessToken)
	assert.Error(t, err)
}