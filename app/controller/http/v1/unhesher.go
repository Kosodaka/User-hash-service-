package v1

import (
	"github.com/gin-gonic/gin"
	unhasherUc "mainHashService/app/usecase"
	"mainHashService/pkg/logger"
	"net/http"
)

type HasherRouter struct {
	logger   logger.Logger
	unhasher unhasherUc.UnhasherUC
}

type HesherRouterParams struct {
	Logger   *logger.Logger
	Unhasher unhasherUc.UnhasherUC
}

func NewHasherRouter(handler *gin.RouterGroup, params *HesherRouterParams) {
	r := &HasherRouter{
		logger:   *params.Logger,
		unhasher: params.Unhasher,
	}

	h := handler.Group("/")
	{
		hasher := h.Group("/")
		{
			hasher.POST("/hash-from-query", r.UnhashFromQuery)
			hasher.POST("/hash-from-file", r.UnhashFromFile)
			hasher.POST("/get-hashed", r.GetHashedData)
		}
		uploader := h.Group("/upload")
		{
			uploader.POST("/", r.UploadFile)
		}
	}
}

func (r *HasherRouter) UnhashFromQuery(ctx *gin.Context) {
	r.logger.Info("call check hash")
	var query UnhasherRequest
	err := ctx.ShouldBindJSON(&query)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url, bucketName, err := r.unhasher.UnhashFromQuery(ctx, query.Query)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url, "bucket": bucketName, "status": "OK"})
}

func (r *HasherRouter) UnhashFromFile(ctx *gin.Context) {
	var bucket UnhashFromFileRequest
	err := ctx.ShouldBindJSON(&bucket)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url, bucketName, err := r.unhasher.UnhashFromFile(ctx, bucket.Bucket, bucket.ObjName)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url, "bucket": bucketName, "status": "OK"})
}

func (r *HasherRouter) UploadFile(ctx *gin.Context) {
	r.logger.Info("call check hash")

	file, err := ctx.FormFile("file")
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileOpened, err := file.Open()
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, bucketName, err := r.unhasher.UplooadFile(ctx, fileOpened, file.Filename, file.Size)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url, "bucket": bucketName, "status": "OK"})
}

func (r *HasherRouter) GetHashedData(ctx *gin.Context) {
	r.logger.Info("call check hash")
	var query UnhasherRequest
	err := ctx.ShouldBindJSON(&query)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.unhasher.GetHashedFile(ctx, query.Query)
	if err != nil {
		r.logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}
