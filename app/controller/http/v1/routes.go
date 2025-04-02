package v1

import (
	"github.com/gin-gonic/gin"
	checkerUc "mainHashService/app/usecase/checker"
	unhasherUc "mainHashService/app/usecase/unhasher"
	"mainHashService/pkg/logger"
)

type RouterParams struct {
	Logger    *logger.Logger
	CheckerUc checkerUc.CheckerUC
	UnesherUc unhasherUc.UnhasherUC
}

func NewRouter(handler *gin.Engine, params *RouterParams) {
	handler.MaxMultipartMemory = 100 << 20
	handler.Use(gin.Recovery())

	h := handler.Group("/api")
	{
		checkerGroup := h.Group("/checker")
		{
			NewCheckerRouter(checkerGroup, &CheckerRouterParams{
				Logger:    params.Logger,
				CheckerUC: params.CheckerUc,
			})
		}
		hasherGroup := h.Group("/hasher")
		{

			NewHasherRouter(hasherGroup, &HesherRouterParams{
				Logger:   params.Logger,
				Unhasher: params.UnesherUc,
			})
		}
	}
}
