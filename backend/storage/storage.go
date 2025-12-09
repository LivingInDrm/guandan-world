package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, key string, data io.Reader, contentType string) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
