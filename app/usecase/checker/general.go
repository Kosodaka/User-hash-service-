package checker

import (
	"context"
	"mainHashService/app/entity"
	"mainHashService/pkg/logger"
)

type CheckerUCImpl struct {
	CheckerRepo CheckerRepo
	lg          logger.Logger
}

func New(lg *logger.Logger, checkerRepo CheckerRepo) *CheckerUCImpl {
	return &CheckerUCImpl{
		CheckerRepo: checkerRepo,
		lg:          *lg,
	}
}

var _ CheckerUC = (*CheckerUCImpl)(nil)

// CheckHash - метод для проверки хеша, сверяет
func (uc *CheckerUCImpl) CheckHash(ctx context.Context, hash CheckerUcDto) (bool, error) {
	hashes := entity.Checker{
		Hash:   hash.Hash,
		Domain: hash.Domain,
	}
	return uc.CheckerRepo.CheckHash(ctx, hashes)
}
