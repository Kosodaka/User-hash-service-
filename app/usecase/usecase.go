package usecase

import (
	"context"
	"io"
)

type UnhasherUC interface {
	UnhashFromQuery(ctx context.Context, query string) (string, string, error)
	UnhashFromFile(ctx context.Context, bucket, objectName string) (string, string, error)
	UplooadFile(ctx context.Context, reader io.ReadCloser, objectName string, size int64) (string, string, error)
	GetHashedFile(ctx context.Context, query string) error
}
