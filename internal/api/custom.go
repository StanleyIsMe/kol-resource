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
	"kolresource/internal/kol/repository/email"
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
	emailRepository := email.NewRepository(a.cfg)
	kolUseCase := kolUseCase.NewKolUseCaseImpl(kolRepository, emailRepository)

	httpRouter.Use(
		middleware.Cors(),
		pkgMiddleware.GinRecover(),
		pkgMiddleware.GinContextLogger(a.logger), //nolint:contextcheck
		pkgMiddleware.GinTimeout(defaultTimeout), //nolint:contextcheck
	)

	adminHTTP.RegisterAdminRoutes(httpRouter, adminUseCase)

	authRouter := httpRouter.Group("", middleware.JWT(adminUseCase))

	kolHTTP.RegisterKolRoutes(authRouter, kolUseCase)
}
