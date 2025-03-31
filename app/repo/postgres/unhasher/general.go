package unhasher

import (
	"mainHashService/app/repo/postgres"
	"mainHashService/pkg/logger"
	"net/http"
)

type RepoImpl struct {
	lg             *logger.Logger
	unhashEndpoint string
	unhashClient   http.Client
}

func New(lg *logger.Logger, unhashEndpoint string) *RepoImpl {
	return &RepoImpl{
		lg:             lg,
		unhashEndpoint: unhashEndpoint,
	}
}

var _ postgres.UnhasherRepo = (*RepoImpl)(nil)
