package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserEntity struct {
	ID           string         `db:"id"`
	Username     string         `db:"username"`
	PasswordHash string         `db:"password_hash"`
	Nickname     sql.NullString `db:"nickname"`
	AvatarKey    sql.NullString `db:"avatar_key"`
	Email        sql.NullString `db:"email"`
	Phone        sql.NullString `db:"phone"`
	WechatOpenID sql.NullString `db:"wechat_openid"`
	Online       bool           `db:"online"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

func (e *UserEntity) ToUser() *User {
	user := &User{
		ID:       e.ID,
		Username: e.Username,
		Password: e.PasswordHash,
		Online:   e.Online,
	}
	if e.Nickname.Valid {
		user.Nickname = &e.Nickname.String
	}
	if e.AvatarKey.Valid {
		user.AvatarKey = &e.AvatarKey.String
	}
	return user
}

type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) (*UserEntity, error)
	FindByUsername(ctx context.Context, username string) (*UserEntity, error)
	FindByID(ctx context.Context, id string) (*UserEntity, error)
	UpdateOnlineStatus(ctx context.Context, id string, online bool) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	UpdateProfile(ctx context.Context, id string, nickname, avatarKey *string) error
}

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, username, passwordHash string) (*UserEntity, error) {
	query := `
		INSERT INTO users (username, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, username, password_hash, nickname, avatar_key, email, phone, wechat_openid, online, created_at, updated_at
	`

	var user UserEntity
	err := r.db.QueryRowxContext(ctx, query, username, passwordHash).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*UserEntity, error) {
	query := `
		SELECT id, username, password_hash, nickname, avatar_key, email, phone, wechat_openid, online, created_at, updated_at 
		FROM users 
		WHERE username = $1
	`

	var user UserEntity
	err := r.db.QueryRowxContext(ctx, query, username).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*UserEntity, error) {
	query := `
		SELECT id, username, password_hash, nickname, avatar_key, email, phone, wechat_openid, online, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	var user UserEntity
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateOnlineStatus(ctx context.Context, id string, online bool) error {
	query := `UPDATE users SET online = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, online, id)
	return err
}

func (r *PostgresUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRepository) UpdateProfile(ctx context.Context, id string, nickname, avatarKey *string) error {
	if nickname == nil && avatarKey == nil {
		return nil
	}
	query := `UPDATE users SET nickname = COALESCE($1, nickname), avatar_key = COALESCE($2, avatar_key), updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, nickname, avatarKey, id)
	return err
}
