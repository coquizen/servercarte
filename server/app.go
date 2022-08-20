package server

import (
	"log"
	"net/http"

	"github.com/coquizen/servercarte/domain/security"

	"github.com/coquizen/servercarte/internal/delivery/ginHTTP"

	"github.com/coquizen/servercarte/internal/authentication/framework/jwt"
	"github.com/coquizen/servercarte/internal/config"
	"github.com/coquizen/servercarte/internal/security/bcrypto"
	"github.com/coquizen/servercarte/internal/store/gormDB"

	"github.com/coquizen/servercarte/domain/account"
	"github.com/coquizen/servercarte/domain/authentication"
	"github.com/coquizen/servercarte/domain/menu"
	"github.com/coquizen/servercarte/domain/user"
	accountTransport "github.com/coquizen/servercarte/internal/account/delivery/ginHTTP"
	accountRepo "github.com/coquizen/servercarte/internal/account/repository/gorm"
	authHTTP "github.com/coquizen/servercarte/internal/authentication/delivery/ginHTTP"
	menuTransport "github.com/coquizen/servercarte/internal/menu/delivery/ginHTTP"
	menuRepo "github.com/coquizen/servercarte/internal/menu/repository/gorm"
	userTransport "github.com/coquizen/servercarte/internal/user/delivery/ginHTTP"
	userRepo "github.com/coquizen/servercarte/internal/user/repository/gorm"
)

// App struct represents this application
type App struct {
	httpServer *http.Server
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
	accountTransport.RegisterRoutes(accountService, authenticationService, ginHandler, authenticationMiddleware, ginHTTP.AuthorizationMiddleware(0))

	server := ginHTTP.NewServer(rCfg, ginHandler)

	return &App{
		server,
	}
}

func (a *App) Run() error {
	return a.httpServer.ListenAndServe()
}
