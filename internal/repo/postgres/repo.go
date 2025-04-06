package postgres

import (
	"context"
	"io"
	"mainHashService/internal/entity"
)

type FetchDataRepo interface {
	GetHashFromQuery(ctx context.Context, query string, args []interface{}) ([]UserData, error)
	GetHashFromFile(ctx context.Context, reader io.ReadCloser) ([]UserData, error)
	GetHashedData(ctx context.Context, query string) ([]HashedData, error)
	QueryBuilder(fields []string, filters []entity.QueryStmt) (string, []interface{}, error)
}
