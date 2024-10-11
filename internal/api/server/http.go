package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetupHTTPServer() {
	ginMode := gin.ReleaseMode
	if s.cfg.Debug {
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)
	s.httpRouter = gin.New()

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.HTTP.Port),
		ReadTimeout:       s.cfg.HTTP.Timeouts.ReadTimeout,
		ReadHeaderTimeout: s.cfg.HTTP.Timeouts.ReadHeaderTimeout,
		WriteTimeout:      s.cfg.HTTP.Timeouts.WriteTimeout,
		IdleTimeout:       s.cfg.HTTP.Timeouts.IdleTimeout,
		Handler:           s.httpRouter,
	}
}

func (s *Server) HTTPRouter() *gin.Engine {
	return s.httpRouter
}

func (s *Server) startHTTPServer() {
	// Start http server
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Err(err).Msg("http server failed to listen and serve")
		}
	}()
}
