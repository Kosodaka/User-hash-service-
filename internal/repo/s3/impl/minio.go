package impl

import (
	s3repo "mainHashService/internal/repo/s3"
	"mainHashService/pkg/logger"
	"mainHashService/pkg/s3"
)

type RepoImpl struct {
	lg logger.Logger
	s3 *s3.S3
}

func New(lg *logger.Logger, s3 *s3.S3) *RepoImpl {
	return &RepoImpl{
		lg: *lg,
		s3: s3,
	}
}

var _ s3repo.ObjsectsRepo = (*RepoImpl)(nil)
