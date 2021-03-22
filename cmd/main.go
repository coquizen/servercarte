package main

import (
	"flag"
	"github.com/CaninoDev/gastro/server/api/account"
	"github.com/CaninoDev/gastro/server/api/menu"
	"github.com/CaninoDev/gastro/server/api/user"
	"github.com/CaninoDev/gastro/server/internal/security/bcrypto"
	"github.com/CaninoDev/gastro/server/internal/storage/gormDB"
	"github.com/CaninoDev/gastro/server/internal/transport/ginHTTP"
	"log"

	"github.com/CaninoDev/gastro/server/internal/authentication/jwt"
	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/logger"
	"github.com/CaninoDev/gastro/server/internal/router/ginRouter"
)

var (
	configYAML       = flag.String("c", "config.yml", "configure db")
	populateDatabase = flag.Bool("p", false, "reset and populate the db with sample data")
)

func main() {
	flag.Parse()
	routerC, databaseC, securityC, authC, err := config.Load(*configYAML)
	if err != nil {
		logger.Error.Fatalf("error parsing config.yml %v", err)
	}

	gDB, err := gormDB.Start(&databaseC, *populateDatabase)
	if err != nil {
		log.Panic(err)
	}

	authService, err := jwt.New(authC)
	if err != nil {
		log.Panic(err)
	}

	passwordService := bcrypto.Initialize(securityC)

	ginHandler := ginRouter.NewGinEngineHandler()

	menuRepository := gormDB.NewMenuRepository(gDB)
	menuService := menu.Initialize(menuRepository)

	userRepository := gormDB.NewUserRepository(gDB)
	userService := user.Initialize(userRepository)

	accountRepository := gormDB.NewAccountRepository(gDB)
	accountService := account.Initialize(accountRepository, userRepository, passwordService, authService)

	ginHTTP.NewMenuRoutes(menuService, authService, ginHandler)
	ginHTTP.NewUserRoutes(userService, ginHandler)
	ginHTTP.NewAccountRoutes(accountService, authService, ginHandler)

	router := ginRouter.Initialize(routerC, ginHandler)

	if err := router.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
