package s3

import (
	"github.com/minio/madmin-go/v2"
	"github.com/minio/minio-go/v7"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	minioCredentials "github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	Aws   *awsS3.S3
	Minio *minio.Client
	Admin *madmin.AdminClient
}

type Config struct {
	Endpoint  string
	AccessID  string
	SecretKey string
	EnableTLS bool
}

func newAws(config *Config) (awsCli *awsS3.S3, err error) {
	var httpClient *http.Client

	awsSession, err := session.NewSession(&aws.Config{
		HTTPClient:       httpClient,
		Region:           aws.String("ru-ru"),
		Endpoint:         aws.String(config.Endpoint),
		DisableSSL:       aws.Bool(!config.EnableTLS),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      awsCredentials.NewStaticCredentials(config.AccessID, config.SecretKey, ""),
	})

	if err != nil {
		return nil, err
	}

	awsCli = awsS3.New(awsSession)

	return awsCli, nil
}

func newMinio(config *Config) (mc *minio.Client, err error) {
	var httpRoundTripper http.RoundTripper

	mc, err = minio.New(config.Endpoint, &minio.Options{
		Transport: httpRoundTripper,
		Region:    "ru-ru",
		Creds:     minioCredentials.NewStaticV4(config.AccessID, config.SecretKey, ""),
		Secure:    false,
	})

	return mc, err
}

func newAdmin(config *Config) (admin *madmin.AdminClient, err error) {
	admin, err = madmin.New(config.Endpoint, config.AccessID, config.SecretKey, config.EnableTLS)

	return admin, err
}

func New(config *Config) (s3 *S3, err error) {
	awsCli, err := newAws(config)
	if err != nil {
		return nil, err
	}

	mc, err := newMinio(config)
	if err != nil {
		return nil, err
	}

	admin, err := newAdmin(config)
	if err != nil {
		return nil, err
	}

	return &S3{
		Aws:   awsCli,
		Minio: mc,
		Admin: admin,
	}, nil
}
