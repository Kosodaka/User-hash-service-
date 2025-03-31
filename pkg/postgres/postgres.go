package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultConnAttempts = 3
	defaultConnTimeout  = time.Second
	defaultMinConns     = 0
	defaultMaxConns     = 16
)

type Config struct {
	PoolConfig pgxpool.Config

	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SslMode  string
	TimeZone string

	ConnAttempts int
	ConnTimeout  time.Duration
	MinConns     int
	MaxConns     int
}

func (cfg *Config) GetDSN() string {
	parts := []string{
		"host=" + cfg.Host,
		"port=" + cfg.Port,
		"user=" + cfg.User,
		"password=" + cfg.Password,
		"dbname=" + cfg.DBName,
		"sslmode=" + cfg.SslMode,
	}

	return strings.Join(parts, " ")
}

func ParseConfig(cfg *Config) (err error) {
	if cfg.ConnAttempts == 0 {
		cfg.ConnAttempts = defaultConnAttempts
	}

	if cfg.ConnTimeout == 0 {
		cfg.ConnTimeout = defaultConnTimeout
	}

	if cfg.MinConns == 0 {
		cfg.MinConns = defaultMinConns
	}

	if cfg.MaxConns == 0 {
		cfg.MaxConns = defaultMaxConns
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return fmt.Errorf("could not parse postgres config: %w", err)
	}

	cfg.PoolConfig = *poolConfig
	cfg.PoolConfig.MinConns = int32(cfg.MinConns)
	cfg.PoolConfig.MaxConns = int32(cfg.MaxConns)

	return nil
}

func NewPgxPool(config *Config) (pool *pgxpool.Pool, err error) {
	for attempts := config.ConnAttempts; attempts > 0; attempts-- {
		pool, err = pgxpool.NewWithConfig(context.Background(), &config.PoolConfig)
		if err == nil {
			break
		}

		log.Printf("trying to connect to postgres, attempts left: %d", attempts)
		time.Sleep(config.ConnTimeout)
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	return pool, nil
}
