package repository

import (
	"context"
	"fmt"

	"kol/internal/api/config"
	"kol/pkg/database/postgreinit"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// NewPGPool to return a connection pool by database config.
func NewPGPool(ctx context.Context, cfg *config.Database, logger *zerolog.Logger) (*pgxpool.Pool, error) {
	pgi, err := postgreinit.New(
		&postgreinit.Config{
			Host:         cfg.Host,
			Port:         cfg.Port,
			User:         cfg.User,
			Password:     cfg.Password,
			Database:     cfg.Database,
			MaxConns:     cfg.MaxConns,
			MaxIdleConns: cfg.MaxIdleConns,
			MaxLifeTime:  cfg.MaxLifeTime,
		},
		postgreinit.WithLogger(logger, "request-id"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PGInit: %w", err)
	}

	pool, err := pgi.ConnPool(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate connection pool: %w", err)
	}

	return pool, nil
}
