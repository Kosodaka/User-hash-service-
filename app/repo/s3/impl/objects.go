package impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"io"
	"time"
)

func (r *RepoImpl) UploadObject(ctx context.Context, filePath, objectName string) (string, string, error) {
	bucketName := generateBucketName()

	exists, err := r.s3.Minio.BucketExists(ctx, bucketName)
	if err != nil {
		return "", "", newErrorFromAws(bucketName, objectName, err)
	}
	if !exists {
		err = r.s3.Minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", "", newErrorFromAws(bucketName, objectName, err)
		}
	}

	_, err = r.s3.Minio.FPutObject(
		ctx,
		bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return "", "", newErrorFromAws(bucketName, objectName, err)
	}

	presignedURL, err := r.s3.Minio.PresignedGetObject(
		ctx,
		bucketName,
		objectName,
		24*time.Hour,
		nil,
	)
	if err != nil {
		return "", "", newErrorFromAws(bucketName, objectName, err)
	}

	return presignedURL.String(), bucketName, nil
}

func generateBucketName() string {
	return uuid.New().String()
}

func (r *RepoImpl) DownloadObject(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	objectReader, err := r.s3.Minio.GetObject(
		ctx,
		bucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, newErrorFromAws(bucket, objectName, err)
	}

	return objectReader, nil
}

func (r *RepoImpl) UploadObjectFromFile(
	ctx context.Context,
	reader io.Reader,
	objectName string,
	size int64,
) (string, string, error) {
	bucketName := generateBucketName()

	exists, err := r.s3.Minio.BucketExists(ctx, bucketName)
	if err != nil {
		return "", "", newErrorFromAws(bucketName, objectName, err)
	}

	if !exists {
		if err := r.s3.Minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return "", "", newErrorFromAws(bucketName, objectName, err)
		}
	}

	opts := minio.PutObjectOptions{
		ContentType: "application/octet-stream",
		PartSize:    64 << 20, // 64MB для multipart upload
	}

	_, err = r.s3.Minio.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		size,
		opts,
	)
	if err != nil {
		return "", "", newErrorFromAws(bucketName, objectName, err)
	}

	presignedURL, err := r.s3.Minio.PresignedGetObject(
		ctx,
		bucketName,
		objectName,
		24*time.Hour,
		nil,
	)
	if err != nil {
		return "", "", newErrorFromAws(bucketName, objectName, err)
	}

	return presignedURL.String(), bucketName, nil
}
