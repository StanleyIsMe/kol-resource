package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"kolresource/internal/api"
	apiCfg "kolresource/internal/api/config"
	"kolresource/pkg/config"
	"kolresource/pkg/shutdown"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const shutdownGracePeriod = 30 * time.Second

// @title           KOL Resource API
// @version         1.0
// @description     API Server for KOL Resource Management System

// @contact.name   Stanley Hsieh
// @contact.email  grimmh6838@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	mainCtx, mainStopCtx := context.WithCancel(context.Background())

	cfg, err := config.LoadWithEnv[apiCfg.Config](mainCtx, "./config/api")
	if err != nil {
		log.Fatal().Err(err).Msg("load config failed")
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if cfg.PrettyLog {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	shutdownHandler := shutdown.New(
		&logger,
		shutdown.WithGracePeriodDuration(shutdownGracePeriod),
	)

	app := api.NewAPI(cfg, shutdownHandler, &logger)
	if err = app.Start(mainCtx); err != nil {
		logger.Fatal().Err(err).Msg("api serve failed with an error")
	}

	if err := shutdownHandler.Listen(
		mainCtx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	); err != nil {
		logger.Fatal().Err(err).Msg("graceful shutdown failed.. forcing exit.")
	}

	mainStopCtx()
}
