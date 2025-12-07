package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Password validation rules
// NOTE: 密码规则需前后端保持一致
// 后端: backend/auth/service.go
// 前端: RegisterForm.tsx, LoginForm.tsx
const (
	MinPasswordLength = 8
	MaxPasswordLength = 50
	SpecialChars      = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func validatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}
	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password must not exceed %d characters", MaxPasswordLength)
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case strings.ContainsRune(SpecialChars, c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// Username validation rules
// NOTE: 用户名规则需前后端保持一致
// 后端: backend/auth/service.go
// 前端: RegisterForm.tsx, LoginForm.tsx
const (
	MinUsernameLength = 4
	MaxUsernameLength = 20
	UsernamePattern   = `^[a-zA-Z0-9_]+$`
)

var usernameRegex = regexp.MustCompile(UsernamePattern)

func validateUsername(username string) error {
	if len(username) < MinUsernameLength {
		return fmt.Errorf("username must be at least %d characters", MinUsernameLength)
	}
	if len(username) > MaxUsernameLength {
		return fmt.Errorf("username must not exceed %d characters", MaxUsernameLength)
	}
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Online   bool   `json:"online"`
}

type AuthToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       string    `json:"user_id"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(username, password string) (*User, error)
	Login(username, password string) (*TokenPair, error)
	ValidateToken(token string) (*User, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	Logout(refreshToken string) error
	GetUserByID(userID string) (*User, error)
}

type authService struct {
	userRepo           UserRepository
	tokenRepo          TokenRepository
	accessSecret       []byte
	refreshSecret      []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type AuthServiceConfig struct {
	AccessSecret       string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func NewAuthService(userRepo UserRepository, tokenRepo TokenRepository, cfg AuthServiceConfig) AuthService {
	return &authService{
		userRepo:           userRepo,
		tokenRepo:          tokenRepo,
		accessSecret:       []byte(cfg.AccessSecret),
		refreshSecret:      []byte(cfg.RefreshSecret),
		accessTokenExpiry:  cfg.AccessTokenExpiry,
		refreshTokenExpiry: cfg.RefreshTokenExpiry,
	}
}

func (s *authService) Register(username, password string) (*User, error) {
	ctx := context.Background()

	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userEntity, err := s.userRepo.Create(ctx, username, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &User{
		ID:       userEntity.ID,
		Username: userEntity.Username,
		Online:   userEntity.Online,
	}, nil
}

func (s *authService) Login(username, password string) (*TokenPair, error) {
	ctx := context.Background()

	userEntity, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if userEntity == nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	tokenPair, err := s.generateTokenPair(userEntity.ID, userEntity.Username)
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(s.refreshTokenExpiry)
	if err := s.tokenRepo.SaveRefreshToken(ctx, userEntity.ID, tokenPair.RefreshToken, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	if err := s.userRepo.UpdateOnlineStatus(ctx, userEntity.ID, true); err != nil {
		return nil, fmt.Errorf("failed to update online status: %w", err)
	}

	return tokenPair, nil
}

func (s *authService) ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	ctx := context.Background()
	userEntity, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if userEntity == nil {
		return nil, errors.New("user not found")
	}

	return userEntity.ToUser(), nil
}

func (s *authService) RefreshToken(refreshToken string) (*TokenPair, error) {
	ctx := context.Background()

	tokenEntity, err := s.tokenRepo.RevokeToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	if tokenEntity == nil {
		return nil, errors.New("invalid or already used refresh token")
	}

	userEntity, err := s.userRepo.FindByID(ctx, tokenEntity.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if userEntity == nil {
		return nil, errors.New("user not found")
	}

	newTokenPair, err := s.generateTokenPair(userEntity.ID, userEntity.Username)
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(s.refreshTokenExpiry)
	if err := s.tokenRepo.SaveRefreshToken(ctx, userEntity.ID, newTokenPair.RefreshToken, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return newTokenPair, nil
}

func (s *authService) Logout(refreshToken string) error {
	ctx := context.Background()

	tokenEntity, err := s.tokenRepo.RevokeToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	if tokenEntity == nil {
		return nil
	}

	if err := s.userRepo.UpdateOnlineStatus(ctx, tokenEntity.UserID, false); err != nil {
		return fmt.Errorf("failed to update online status: %w", err)
	}

	return nil
}

func (s *authService) GetUserByID(userID string) (*User, error) {
	ctx := context.Background()

	userEntity, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if userEntity == nil {
		return nil, errors.New("user not found")
	}

	return userEntity.ToUser(), nil
}

func (s *authService) generateTokenPair(userID, username string) (*TokenPair, error) {
	accessExpiresAt := time.Now().Add(s.accessTokenExpiry)
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenString := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiresAt,
		UserID:       userID,
	}, nil
}
