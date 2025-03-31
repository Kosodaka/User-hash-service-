package unhasher

import (
	"mainHashService/app/repo/postgres"
	s3repo "mainHashService/app/repo/s3"
)

type UnhasherRepo interface {
	postgres.UnhasherRepo
}

type FetchRepo interface {
	postgres.FetchDataRepo
}

type S3Repo interface {
	s3repo.ObjsectsRepo
}
