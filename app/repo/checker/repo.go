package checker

import (
	"context"
	"mainHashService/app/entity"
)

type CheckerRepo interface {
	CheckHash(ctx context.Context, hash entity.Checker) (bool, error)
}
