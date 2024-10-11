package api

import (
	"context"
	"database/sql"
	"fmt"

	"kolresource/internal/admin/delivery/http"
	adminRepo "kolresource/internal/admin/repository"
	adminUseCase "kolresource/internal/admin/usecase"
	apiCfg "kolresource/internal/api/config"
	"kolresource/internal/api/repository"
	"kolresource/internal/api/server"
	"kolresource/pkg/config"
	"kolresource/pkg/shutdown"

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

	pgStdConn, err := repository.NewPGStdConn(ctx, &a.cfg.CustomConfig.DB, a.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize PGStdConn: %w", err)
	}

	a.shutdownHandler.Add("pgStdConn", func(ctx context.Context) error {
		return pgStdConn.Close()
	})

	a.registerHTTPSvc(ctx, pgStdConn)

	if err := apiS.Start(ctx); err != nil {
		return fmt.Errorf("server start failed: %w", err)
	}

	a.shutdownHandler.Add("server", apiS.Shutdown)

	return nil
}

func (a *API) registerHTTPSvc(_ context.Context, dbStdConn *sql.DB) {
	a.server.SetupHTTPServer()
	httpRouter := a.server.HTTPRouter()

	adminRepository := adminRepo.NewAdminRepository(dbStdConn)

	adminUseCase := adminUseCase.NewAdminUseCaseImpl(adminRepository)

	http.RegisterAdminRoutes(httpRouter, adminUseCase)
}
