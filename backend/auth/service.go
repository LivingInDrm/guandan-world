package auth

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Never expose password in JSON
	Online   bool   `json:"online"`
}

// AuthToken represents an authentication token
type AuthToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
}

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthService interface defines authentication operations
type AuthService interface {
	Register(username, password string) (*User, error)
	Login(username, password string) (*AuthToken, error)
	ValidateToken(token string) (*User, error)
	Logout(token string) error
	GetUserByID(userID string) (*User, error)
}

// authService implements AuthService interface
type authService struct {
	users        map[string]*User // userID -> User
	usersByName  map[string]*User // username -> User
	tokens       map[string]*AuthToken // token -> AuthToken
	jwtSecret    []byte
	tokenExpiry  time.Duration
	mu           sync.RWMutex
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecret string, tokenExpiry time.Duration) AuthService {
	return &authService{
		users:       make(map[string]*User),
		usersByName: make(map[string]*User),
		tokens:      make(map[string]*AuthToken),
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

// Register creates a new user account
func (s *authService) Register(username, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate input
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Check if username already exists
	if _, exists := s.usersByName[username]; exists {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate user ID
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())

	// Create user
	user := &User{
		ID:       userID,
		Username: username,
		Password: string(hashedPassword),
		Online:   false,
	}

	// Store user
	s.users[userID] = user
	s.usersByName[username] = user

	// Return user without password
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Online:   user.Online,
	}, nil
}

// Login authenticates a user and returns a token
func (s *authService) Login(username, password string) (*AuthToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find user by username
	user, exists := s.usersByName[username]
	if !exists {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	expiresAt := time.Now().Add(s.tokenExpiry)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create auth token
	authToken := &AuthToken{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
	}

	// Store token
	s.tokens[tokenString] = authToken

	// Mark user as online
	user.Online = true

	return authToken, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *authService) ValidateToken(tokenString string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if token exists in our store
	authToken, exists := s.tokens[tokenString]
	if !exists {
		return nil, errors.New("invalid token")
	}

	// Check if token is expired
	if time.Now().After(authToken.ExpiresAt) {
		// Clean up expired token
		delete(s.tokens, tokenString)
		return nil, errors.New("token expired")
	}

	// Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Get user
	user, exists := s.users[claims.UserID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Logout invalidates a token
func (s *authService) Logout(tokenString string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find and remove token
	authToken, exists := s.tokens[tokenString]
	if !exists {
		return errors.New("token not found")
	}

	// Mark user as offline
	if user, exists := s.users[authToken.UserID]; exists {
		user.Online = false
	}

	// Remove token
	delete(s.tokens, tokenString)

	return nil
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(userID string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}