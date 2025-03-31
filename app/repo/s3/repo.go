package s3

import (
	"context"
	"io"
)

type ObjsectsRepo interface {
	UploadObject(ctx context.Context, filePath, objectName string) (string, string, error)
	DownloadObject(ctx context.Context, bucket, objectName string) (io.ReadCloser, error)
	UploadObjectFromFile(ctx context.Context, reader io.Reader,
		objetName string, size int64) (string, string, error)
}
