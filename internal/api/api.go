package api

import (
	"context"
	"fmt"

	apiCfg "kol/internal/api/config"
	"kol/internal/api/server"
	"kol/pkg/config"
	"kol/pkg/shutdown"

	"github.com/rs/zerolog"
)

type API struct {
	logger          *zerolog.Logger
	cfg             *config.Config[apiCfg.Config]
	server          *server.Server
	shutdownHandler *shutdown.Shutdown
}

// NewAPI to return an API instance to support Serve/Shutdown
func NewAPI(cfg *config.Config[apiCfg.Config], shutdownHandler *shutdown.Shutdown, logger *zerolog.Logger) *API {
	return &API{
		logger:          logger,
		cfg:             cfg,
		shutdownHandler: shutdownHandler,
	}
}

func (a *API) Start(ctx context.Context) error {
	// api server
	apiS := server.NewServer(a.cfg, a.logger)
	a.server = apiS

	a.registerHTTPSvc(ctx)

	if err := apiS.Start(ctx); err != nil {
		return fmt.Errorf("server start failed: %w", err)
	}

	a.shutdownHandler.Add("server", apiS.Shutdown)

	return nil
}

func (a *API) registerHTTPSvc(_ context.Context) {
	a.server.SetupHTTPServer()

	//
}