package api

import (
	"context"
	"database/sql"
	"time"

	adminHTTP "kolresource/internal/admin/delivery/http"
	adminRepo "kolresource/internal/admin/repository"
	adminUseCase "kolresource/internal/admin/usecase"
	"kolresource/internal/api/middleware"
	emailHTTP "kolresource/internal/email/delivery/http"
	EmailRepo "kolresource/internal/email/repository/sqlboiler"
	emailUseCase "kolresource/internal/email/usecase"
	kolHTTP "kolresource/internal/kol/delivery/http"
	kolRepo "kolresource/internal/kol/repository/sqlboiler"
	kolUseCase "kolresource/internal/kol/usecase"
	pkgMiddleware "kolresource/pkg/transport/middleware"
)

const (
	defaultTimeout = 10 * time.Second
)

func (a *API) registerHTTPSvc(ctx context.Context, dbStdConn *sql.DB) {
	a.server.SetupHTTPServer()
	httpRouter := a.server.HTTPRouter()

	adminRepository := adminRepo.NewAdminRepository(dbStdConn)

	adminUseCase := adminUseCase.NewAdminUseCaseImpl(adminRepository, a.cfg)

	kolRepository := kolRepo.NewKolRepository(dbStdConn)
	emailRepository := EmailRepo.NewEmailRepository(dbStdConn)

	kolUseCase := kolUseCase.NewKolUseCaseImpl(kolRepository)
	emailUseCase := emailUseCase.NewEmailUseCaseImpl(emailRepository, kolUseCase)

	httpRouter.Use(
		middleware.Cors(),
		pkgMiddleware.GinRecover(),
		pkgMiddleware.GinContextLogger(a.logger), //nolint:contextcheck
		pkgMiddleware.GinTimeout(defaultTimeout), //nolint:contextcheck
	)

	// admin domain
	adminHTTP.RegisterAdminRoutes(httpRouter, adminUseCase)

	authRouter := httpRouter.Group("", middleware.JWT(adminUseCase))

	// kol domain
	kolHTTP.RegisterKolRoutes(authRouter, kolUseCase)

	// email domain
	emailHTTP.RegisterEmailRoutes(ctx, authRouter, emailHTTP.RegisterEmailRoutesParams{
		Router:          authRouter,
		Cfg:             a.cfg,
		EmailUsecase:    emailUseCase,
		EmailRepository: emailRepository,
	})
}
