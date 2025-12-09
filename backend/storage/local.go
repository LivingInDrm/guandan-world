package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) validatePath(key string) (string, error) {
	cleanKey := filepath.Clean(key)
	if strings.Contains(cleanKey, "..") {
		return "", fmt.Errorf("invalid key: path traversal detected")
	}

	fullPath := filepath.Join(s.basePath, cleanKey)
	cleanBase := filepath.Clean(s.basePath) + string(os.PathSeparator)
	if !strings.HasPrefix(fullPath, cleanBase) && fullPath != filepath.Clean(s.basePath) {
		return "", fmt.Errorf("invalid key: path outside base directory")
	}
	return fullPath, nil
}

func (s *LocalStorage) Save(ctx context.Context, key string, data io.Reader, contentType string) error {
	fullPath, err := s.validatePath(key)
	if err != nil {
		return err
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		os.Remove(fullPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath, err := s.validatePath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath, err := s.validatePath(key)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check file: %w", err)
}
