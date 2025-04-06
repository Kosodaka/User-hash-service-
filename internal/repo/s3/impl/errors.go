package impl

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
)

type NoSuchBucketError struct {
	Bucket string
}

func newNoSuchBucketError(bucket string) error {
	return &NoSuchBucketError{
		Bucket: bucket,
	}
}

func (e *NoSuchBucketError) Error() string {
	return "no such bucket \"" + e.Bucket + "\""
}

type NoSuchObjectError struct {
	Bucket string
	Key    string
}

func newNoSuchObjectError(bucket, key string) error {
	return &NoSuchObjectError{
		Bucket: bucket,
		Key:    key,
	}
}

func (e *NoSuchObjectError) Error() string {
	return "no such object \"" + e.Key + "\" in bucket \"" + e.Bucket + "\""
}

func newErrorFromAws(bucket, key string, err error) error {
	var aerr awserr.Error

	if errors.As(err, &aerr) {
		switch aerr.Code() {
		case awsS3.ErrCodeNoSuchBucket:
			return newNoSuchBucketError(bucket)
		case awsS3.ErrCodeNoSuchKey:
			return newNoSuchObjectError(bucket, key)
		}
	}

	return err
}
