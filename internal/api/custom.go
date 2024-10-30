package api

import (
	"context"
	"database/sql"
	"time"

	adminHTTP "kolresource/internal/admin/delivery/http"
	adminRepo "kolresource/internal/admin/repository"
	adminUseCase "kolresource/internal/admin/usecase"
	"kolresource/internal/api/middleware"
	kolHTTP "kolresource/internal/kol/delivery/http"
	emailRepo "kolresource/internal/kol/repository/email"
	kolRepo "kolresource/internal/kol/repository/sqlboiler"
	kolUseCase "kolresource/internal/kol/usecase"
	pkgMiddleware "kolresource/pkg/transport/middleware"
)

const (
	defaultTimeout = 10 * time.Second
)

func (a *API) registerHTTPSvc(_ context.Context, dbStdConn *sql.DB) {
	a.server.SetupHTTPServer()
	httpRouter := a.server.HTTPRouter()

	adminRepository := adminRepo.NewAdminRepository(dbStdConn)

	adminUseCase := adminUseCase.NewAdminUseCaseImpl(adminRepository, a.cfg)

	kolRepository := kolRepo.NewKolRepository(dbStdConn)
	emailRepository := emailRepo.NewEmailRepository(a.cfg.CustomConfig.Email.ServerHost, a.cfg.CustomConfig.Email.ServerPort)
	kolUseCase := kolUseCase.NewKolUseCaseImpl(kolRepository, emailRepository, a.cfg)

	httpRouter.Use(
		middleware.Cors(),
		pkgMiddleware.GinRecover(a.logger),
		pkgMiddleware.GinContextLogger(a.logger),
		pkgMiddleware.GinTimeout(a.logger, defaultTimeout),
	)

	adminHTTP.RegisterAdminRoutes(httpRouter, adminUseCase)

	authRouter := httpRouter.Group("", middleware.JWT(adminUseCase))

	kolHTTP.RegisterKolRoutes(authRouter, kolUseCase)
}
