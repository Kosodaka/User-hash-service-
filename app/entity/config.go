package entity

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"mainHashService/pkg/logger"
	"mainHashService/pkg/postgres"
	"mainHashService/pkg/s3"
)

type Config struct {
	HTTP
	Database
	Minio
	Log
	HashConf
	App
}

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type App struct {
	TempDir string `env:"TEMP_DIR" env-required:"-"`
	ZipPass string `env:"ZIP_PASS" env-required:"-"`
}

type Log struct {
	Level logger.Level `env-required:""  env:"LOG_LEVEL"`
}

type HashConf struct {
	HMACSecert string `env-required:"" env:"HMAC_SECRET"`
	Endpoint   string `env-required:"" env:"UNHASH_ENDPOINT"`
}

type HTTP struct {
	Host string `env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port string `env:"HTTP_PORT" env-default:"3000"`
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASS"`
	DBName   string `env:"DB_NAME"`
	SslMode  string `env:"DB_SSL_MODE"`
	MinConns int    `env:"DB_MIN_CONNS"`
	MaxConns int    `env:"DB_MAX_CONNS"`
}

type Minio struct {
	EndPoint        string ` env:"S3_ENDPOINT"`
	AccessID        string ` env:"S3_ACCESSID"`
	SecretKey       string ` env:"S3_SECRETKEY"`
	EnableTLS       bool   `env-default:"false"  env:"S3_ENABLE_SSL"`
	DownloadSignKey string ` env:"DOWNLOAD_SIGN_KEY"`
	Partition       string ` env:"PARTITION" `
	Service         string `env:"SERVICE" `
	Region          string `env:"REGION" `
	AccountID       string `env:"ACCOUNT_ID" `
	Resource        string ` env:"RESOURCE" `
}

func (cfg *Config) GetConfigForDB() (dbConf postgres.Config) {
	return postgres.Config{
		Host:     cfg.Database.Host,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		Port:     cfg.Database.Port,
		SslMode:  cfg.Database.SslMode,
		MinConns: cfg.Database.MinConns,
		MaxConns: cfg.Database.MaxConns,
	}
}

func (cfg *Config) GetConfigForS3() (s3Cfg s3.Config) {
	return s3.Config{
		Endpoint:  cfg.Minio.EndPoint,
		AccessID:  cfg.Minio.AccessID,
		SecretKey: cfg.Minio.SecretKey,
		EnableTLS: cfg.Minio.EnableTLS,
	}
}
