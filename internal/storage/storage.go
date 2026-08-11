package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, src io.Reader, objectPath string) (string, error)
	Delete(ctx context.Context, objectPath string) error
}
