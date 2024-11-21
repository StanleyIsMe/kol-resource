package server

import (
	"context"
	"fmt"
	apiCfg "kolresource/internal/api/config"
	"kolresource/pkg/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var (
	CommitTime = "dev"
	CommitHash = "dev"
)

type Server struct {
	logger     *zerolog.Logger
	cfg        *config.Config[apiCfg.Config]
	httpServer *http.Server
	httpRouter *gin.Engine
}

func NewServer(cfg *config.Config[apiCfg.Config], logger *zerolog.Logger) *Server {
	srv := &Server{
		cfg:    cfg,
		logger: logger,
	}

	return srv
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info().
		Str("commitTime", CommitTime).
		Str("commitHash", CommitHash).
		Msg(s.cfg.Name)

	s.startHTTPServer()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown with err: %w", err)
	}

	return nil
}
