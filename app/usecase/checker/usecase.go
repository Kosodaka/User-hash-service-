package checker

import (
	"context"
)

type CheckerUC interface {
	CheckHash(ctx context.Context, hash CheckerUcDto) (bool, error)
}
