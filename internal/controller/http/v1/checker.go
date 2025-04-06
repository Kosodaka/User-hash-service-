package v1

import (
	"github.com/gin-gonic/gin"
	"mainHashService/internal/entity"
	checkerUc "mainHashService/internal/usecase"
	"mainHashService/pkg/logger"
)

type CheckerRouter struct {
	logger    logger.Logger
	checkerUC checkerUc.CheckerUC
}

type CheckerRouterParams struct {
	Logger    *logger.Logger
	CheckerUC checkerUc.CheckerUC
}

func NewCheckerRouter(handler *gin.RouterGroup, params *CheckerRouterParams) {

	r := &CheckerRouter{
		logger:    *params.Logger,
		checkerUC: params.CheckerUC,
	}
	h := handler.Group("/")
	{
		checkerGroup := h.Group("/")
		{
			checkerGroup.POST("/hash", r.CheckHash)
		}
	}
}

func (r *CheckerRouter) CheckHash(ctx *gin.Context) {
	r.logger.Info("call check hash")
	var req CheckerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}

	flag, err := r.checkerUC.CheckHash(ctx, entity.Checker{
		Hash:   req.Hash,
		Domain: req.Domain,
	})
	if err != nil {
		r.logger.Error(err.Error())
		ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "OK", "flag": flag})

}
