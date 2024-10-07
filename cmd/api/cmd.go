package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"kol/internal/api"
	apiCfg "kol/internal/api/config"
	"kol/pkg/config"
	"kol/pkg/shutdown"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const shutdownGracePeriod = 30 * time.Second

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
