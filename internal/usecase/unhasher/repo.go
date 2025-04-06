package unhasher

import (
	"mainHashService/internal/repo/postgres"
	s3repo "mainHashService/internal/repo/s3"
)

type FetchRepo interface {
	postgres.FetchDataRepo
}

type S3Repo interface {
	s3repo.ObjsectsRepo
}
