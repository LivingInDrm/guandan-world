package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type RefreshTokenEntity struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	Revoked   bool      `db:"revoked"`
}

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshTokenEntity, error)
	RevokeToken(ctx context.Context, token string) (*RefreshTokenEntity, error)
	RevokeAllUserTokens(ctx context.Context, userID string) error
	CleanExpired(ctx context.Context) (int64, error)
}

type PostgresTokenRepository struct {
	db *sqlx.DB
}

func NewPostgresTokenRepository(db *sqlx.DB) *PostgresTokenRepository {
	return &PostgresTokenRepository{db: db}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *PostgresTokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at) 
		VALUES ($1, $2, $3)
	`
	
	tokenHash := hashToken(token)
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *PostgresTokenRepository) FindRefreshToken(ctx context.Context, token string) (*RefreshTokenEntity, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked 
		FROM refresh_tokens 
		WHERE token_hash = $1 AND revoked = FALSE AND expires_at > NOW()
	`
	
	tokenHash := hashToken(token)
	var entity RefreshTokenEntity
	err := r.db.QueryRowxContext(ctx, query, tokenHash).StructScan(&entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

func (r *PostgresTokenRepository) RevokeToken(ctx context.Context, token string) (*RefreshTokenEntity, error) {
	query := `
		UPDATE refresh_tokens 
		SET revoked = TRUE 
		WHERE token_hash = $1 AND revoked = FALSE AND expires_at > NOW()
		RETURNING id, user_id, token_hash, expires_at, created_at, revoked
	`
	tokenHash := hashToken(token)
	var entity RefreshTokenEntity
	err := r.db.QueryRowxContext(ctx, query, tokenHash).StructScan(&entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *PostgresTokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *PostgresTokenRepository) CleanExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW() OR revoked = TRUE`
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
