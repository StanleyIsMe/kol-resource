package postgreinit

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	zerologadapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

const (
	defaultMaxConns     = 25
	defaultMaxIdleConns = 25
	defaultMaxLifeTime  = 5 * time.Minute
)

// Option configures PostgreInit behaviour.
type Option func(*PostgreInit)

// PostgreInit provides capabilities for connect to postgres with pgx.pool.
type PostgreInit struct {
	pgxConf *pgxpool.Config
	logLvl  tracelog.LogLevel
}

func New(conf *Config, opts ...Option) (*PostgreInit, error) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		conf.User, conf.Password, net.JoinHostPort(conf.Host, conf.Port), conf.Database,
	)

	pgxConf, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pgxConf.MaxConns = defaultMaxConns
	if conf.MaxConns != 0 {
		pgxConf.MaxConns = conf.MaxConns
	}

	pgxConf.MinConns = defaultMaxIdleConns
	if conf.MaxIdleConns != 0 && conf.MaxConns >= conf.MaxIdleConns {
		pgxConf.MinConns = conf.MaxIdleConns
	} else {
		pgxConf.MinConns = pgxConf.MaxConns
	}

	pgxConf.MaxConnLifetime = defaultMaxLifeTime
	if conf.MaxLifeTime != 0 {
		pgxConf.MaxConnLifetime = conf.MaxLifeTime
	}

	pgi := &PostgreInit{
		pgxConf: pgxConf,
		logLvl:  tracelog.LogLevelWarn,
	}

	for _, opt := range opts {
		opt(pgi)
	}

	pgi.pgxConf.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())

		return nil
	}

	return pgi, nil
}

// ConnPool initiates connection to database and return a pgxpool.Pool.
func (pgi *PostgreInit) ConnPool(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, pgi.pgxConf)
	if err != nil {
		return nil, fmt.Errorf("connect config: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func (pgi *PostgreInit) StdConn(ctx context.Context) (*sql.DB, error) {
	connStr := stdlib.RegisterConnConfig(pgi.pgxConf.ConnConfig)

	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("open stdlib connection with config: %w", err)
	}

	dbConn.SetConnMaxIdleTime(pgi.pgxConf.MaxConnIdleTime)
	dbConn.SetConnMaxLifetime(pgi.pgxConf.MaxConnLifetime)
	dbConn.SetMaxOpenConns(int(pgi.pgxConf.MaxConns))
	dbConn.SetMaxIdleConns(int(pgi.pgxConf.MinConns))

	if err := dbConn.PingContext(ctx); err != nil {
		dbConn.Close()

		return nil, fmt.Errorf("ping std connection with config: %w", err)
	}

	return dbConn, nil
}

// WithLogger Add logger to pgx. if the request context contains request id,
// can pass in the request id context key to reqIDKeyFromCtx and logger will
// log with the request id. Only will log if the log level is equal and above pgx.LogLevelWarn.
func WithLogger(logger *zerolog.Logger, reqIDKeyFromCtx string) Option {
	return func(pgi *PostgreInit) {
		zeroLogger := zerologadapter.NewLogger(*logger, zerologadapter.WithContextFunc(
			func(ctx context.Context, logWith zerolog.Context) zerolog.Context {
				if ctxValue, ok := ctx.Value(reqIDKeyFromCtx).(string); ok {
					logWith = logWith.Str(reqIDKeyFromCtx, ctxValue)
				}

				return logWith
			},
		), zerologadapter.WithoutPGXModule())

		tracerLogger := &tracelog.TraceLog{
			Logger:   zeroLogger,
			LogLevel: pgi.logLvl,
		}

		pgi.pgxConf.ConnConfig.Tracer = tracerLogger
	}
}

// WithLogLevel set pgx log level.
func WithLogLevel(zLvl zerolog.Level) Option {
	return func(pgi *PostgreInit) {
		switch zLvl {
		case zerolog.DebugLevel:
			pgi.logLvl = tracelog.LogLevelDebug
		case zerolog.InfoLevel:
			pgi.logLvl = tracelog.LogLevelInfo
		case zerolog.WarnLevel:
			pgi.logLvl = tracelog.LogLevelWarn
		case zerolog.ErrorLevel:
			pgi.logLvl = tracelog.LogLevelError
		case zerolog.NoLevel:
			pgi.logLvl = tracelog.LogLevelNone
		}
	}
}
