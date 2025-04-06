package v1

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"mainHashService/internal/usecase"
	"mainHashService/pkg/logger"
)

type RouterParams struct {
	Logger    *logger.Logger
	CheckerUc usecase.CheckerUC
	UnesherUc usecase.UnhasherUC
}

func NewRouter(handler *gin.Engine, params *RouterParams) {
	handler.MaxMultipartMemory = 100 << 20
	handler.Use(gin.Recovery())
	pprof.Register(handler)
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
