package auth

import (
	"context"
	"errors"
	"sync"
	"time"
)

type InMemoryUserRepository struct {
	users       map[string]*UserEntity
	usersByName map[string]*UserEntity
	mu          sync.RWMutex
	idCounter   int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:       make(map[string]*UserEntity),
		usersByName: make(map[string]*UserEntity),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, username, passwordHash string) (*UserEntity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usersByName[username]; exists {
		return nil, errors.New("username already exists")
	}

	r.idCounter++
	id := generateTestID(r.idCounter)

	user := &UserEntity{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		Online:       false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r.users[id] = user
	r.usersByName[username] = user

	return user, nil
}

func (r *InMemoryUserRepository) FindByUsername(ctx context.Context, username string) (*UserEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByName[username]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*UserEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *InMemoryUserRepository) UpdateOnlineStatus(ctx context.Context, id string, online bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}
	user.Online = online
	user.UpdatedAt = time.Now()
	return nil
}

func (r *InMemoryUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.usersByName[username]
	return exists, nil
}

func (r *InMemoryUserRepository) UpdateProfile(ctx context.Context, id string, nickname, avatarKey *string) error {
	if nickname == nil && avatarKey == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}
	if nickname != nil {
		user.Nickname.Valid = true
		user.Nickname.String = *nickname
	}
	if avatarKey != nil {
		user.AvatarKey.Valid = true
		user.AvatarKey.String = *avatarKey
	}
	user.UpdatedAt = time.Now()
	return nil
}

type InMemoryTokenRepository struct {
	tokens map[string]*RefreshTokenEntity
	mu     sync.RWMutex
}

func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		tokens: make(map[string]*RefreshTokenEntity),
	}
}

func (r *InMemoryTokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tokenHash := hashToken(token)
	r.tokens[tokenHash] = &RefreshTokenEntity{
		ID:        generateTestID(len(r.tokens) + 1),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	}
	return nil
}

func (r *InMemoryTokenRepository) FindRefreshToken(ctx context.Context, token string) (*RefreshTokenEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tokenHash := hashToken(token)
	entity, exists := r.tokens[tokenHash]
	if !exists || entity.Revoked || time.Now().After(entity.ExpiresAt) {
		return nil, nil
	}
	return entity, nil
}

func (r *InMemoryTokenRepository) RevokeToken(ctx context.Context, token string) (*RefreshTokenEntity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tokenHash := hashToken(token)
	entity, exists := r.tokens[tokenHash]
	if !exists || entity.Revoked || time.Now().After(entity.ExpiresAt) {
		return nil, nil
	}
	entity.Revoked = true
	return entity, nil
}

func (r *InMemoryTokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entity := range r.tokens {
		if entity.UserID == userID {
			entity.Revoked = true
		}
	}
	return nil
}

func (r *InMemoryTokenRepository) CleanExpired(ctx context.Context) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int64
	for hash, entity := range r.tokens {
		if entity.Revoked || time.Now().After(entity.ExpiresAt) {
			delete(r.tokens, hash)
			count++
		}
	}
	return count, nil
}

func generateTestID(n int) string {
	return "test-id-" + string(rune('0'+n%10))
}

func NewTestAuthService(accessSecret string, accessExpiry time.Duration) AuthService {
	userRepo := NewInMemoryUserRepository()
	tokenRepo := NewInMemoryTokenRepository()
	return NewAuthService(userRepo, tokenRepo, AuthServiceConfig{
		AccessSecret:       accessSecret,
		RefreshSecret:      accessSecret + "-refresh",
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: accessExpiry * 24,
	})
}
