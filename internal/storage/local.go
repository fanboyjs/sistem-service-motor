package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	dir string
}

func NewLocalStorage(dir string) *LocalStorage {
	return &LocalStorage{dir: dir}
}

func (s *LocalStorage) Save(ctx context.Context, src io.Reader, objectPath string) (string, error) {
	objectPath = strings.TrimPrefix(objectPath, "/")

	// buat folder tujuan jika belum ada
	// 0o755 -> permission `rwxr-xr-x` (owner bisa baca/tulis/eksekusi; group & other hanya baca/eksekusi)
	fullPath := filepath.Join(s.dir, filepath.FromSlash(objectPath))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return "/uploads/" + objectPath, nil
}

func (s *LocalStorage) Delete(ctx context.Context, objectPath string) error {
	objectPath = strings.TrimPrefix(objectPath, "/uploads/")

	fullPath := filepath.Join(s.dir, filepath.FromSlash(objectPath))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
