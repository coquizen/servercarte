package server

import (
	"log"
	"net/http"

	"github.com/CaninoDev/gastro/server/internal/delivery/ginHTTP"

	"github.com/CaninoDev/gastro/server/internal/authentication/framework/jwt"
	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/security/bcrypto"
	"github.com/CaninoDev/gastro/server/internal/store/gormDB"

	"github.com/CaninoDev/gastro/server/authentication"
	"github.com/CaninoDev/gastro/server/domain/account"
	"github.com/CaninoDev/gastro/server/domain/menu"
	"github.com/CaninoDev/gastro/server/domain/user"
	"github.com/CaninoDev/gastro/server/security"

	accountTransport "github.com/CaninoDev/gastro/server/internal/account/delivery/ginHTTP"
	accountRepo "github.com/CaninoDev/gastro/server/internal/account/repository/gorm"
	authHTTP "github.com/CaninoDev/gastro/server/internal/authentication/delivery/ginHTTP"
	menuTransport "github.com/CaninoDev/gastro/server/internal/menu/delivery/ginHTTP"
	menuRepo "github.com/CaninoDev/gastro/server/internal/menu/repository/gorm"
	userTransport "github.com/CaninoDev/gastro/server/internal/user/delivery/ginHTTP"
	userRepo "github.com/CaninoDev/gastro/server/internal/user/repository/gorm"
)

// App struct represents this application
type App struct {
	httpServer     *http.Server
}

// NewApp serves as the main entry point for this application
func NewApp(rCfg config.Router, dCfg config.Database, aCfg config.Authentication, sCfg config.Security,
	seedDatabase bool) *App {
	//Set up repositories
	db, err := gormDB.Start(dCfg, seedDatabase)
	if err != nil {
		log.Panicf("failed loading database: %v", err)
	}

	menuRepository := menuRepo.NewMenuRepository(db)
	userRepository := userRepo.NewUserRepository(db)
	accountRepository := accountRepo.NewAccountRepository(db)

	authenticationFramework, err := jwt.New(aCfg)
	if err != nil {
		log.Panicf("authentication framework loading error %v", err)
	}
	authenticationService := authentication.NewService(authenticationFramework)

	securityFramework := bcrypto.NewSecurityFramework(sCfg)
	securityService := security.NewService(securityFramework)

	// Setup services
	menuService := menu.NewService(menuRepository)
	userService := user.NewService(userRepository)
	accountService := account.NewService(accountRepository, userService, securityService, authenticationService)


	authenticationMiddleware := authHTTP.NewMiddleWare(authenticationService)

	ginHandler := ginHTTP.NewHandler(rCfg)
	menuTransport.RegisterRoutes(menuService, authenticationService, ginHandler, authenticationMiddleware, ginHTTP.AuthorizationMiddleware(0))
	userTransport.RegisterRoutes(userService, ginHandler)
	accountTransport.RegisterRoutes(authenticationService, accountService,
		ginHandler, authenticationMiddleware,ginHTTP.AuthorizationMiddleware(0))

	server := ginHTTP.NewServer(rCfg, ginHandler)

	return &App{
		server,
	}
}

func (a *App) Run() error {
	return a.httpServer.ListenAndServe()
}
