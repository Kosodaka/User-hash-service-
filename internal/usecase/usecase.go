package usecase

import (
	"context"
	"io"
	"mainHashService/internal/entity"
)

type UnhasherUC interface {
	UnhashFromQuery(ctx context.Context, fields []string, filters []entity.QueryStmt) (string, string, error)
	UnhashFromFile(ctx context.Context, bucket, objectName string) (string, string, error)
	UplooadFile(ctx context.Context, reader io.ReadCloser, objectName string, size int64) (string, string, error)
	GetHashedFile(ctx context.Context, query string) error
}

type CheckerUC interface {
	CheckHash(ctx context.Context, hash entity.Checker) (bool, error)
}
