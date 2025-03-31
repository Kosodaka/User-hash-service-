package postgres

import (
	"context"
	"io"
	"mainHashService/app/entity"
)

type UnhasherRepo interface {
	UnhashData(ctx context.Context, data *Unhashdata) (entity.UnhashedData, error)
}

type FetchDataRepo interface {
	GetHashFromQuery(ctx context.Context, query string) ([]UserData, error)
	GetHashFromFile(ctx context.Context, reader io.ReadCloser) ([]UserData, error)
	GetHashedData(ctx context.Context, query string) ([]HashedData, error)
}
