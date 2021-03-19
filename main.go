package main

import (
	"flag"
	"log"

	"github.com/CaninoDev/gastro/server/internal/api/account"
	"github.com/CaninoDev/gastro/server/internal/api/menu"
	"github.com/CaninoDev/gastro/server/internal/api/user"
	"github.com/CaninoDev/gastro/server/internal/authentication/jwt"
	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/db/gormDB"
	"github.com/CaninoDev/gastro/server/internal/logger"
	"github.com/CaninoDev/gastro/server/internal/router/ginRouter"
	"github.com/CaninoDev/gastro/server/internal/security"
)

var (
	configYAML       = flag.String("c", "config.yml", "configure db")
	populateDatabase = flag.Bool("p", false, "reset and populate the db with sample data")
)

func main() {
	flag.Parse()
	routerC, databaseC, securityC, jwtC, err := config.Load(*configYAML)
	if err != nil {
		logger.Error.Fatalf("error parsing config.yml %v", err)
	}

	gDB, err := gormDB.Start(databaseC, *populateDatabase)
	if err != nil {
		log.Panic(err)
	}

	authService, err := jwt.New(jwtC)
	if err != nil {
		log.Panic(err)
	}


	passwordService := security.Initialize(securityC)

	ginHandler := ginRouter.NewGinEngineHandler()

	menuRepository := menu.NewGormDBRepository(gDB)
	menuService := menu.Initialize(menuRepository)

	menu.NewGinRoutes(menuService, authService, ginHandler)

	userRepository := user.NewGormDBRepository(gDB)
	userService := user.Initialize(userRepository)
	user.NewGinRoutes(userService, ginHandler)

	accountRepository := account.NewGormDBRepository(gDB)

	accountService := account.Initialize(accountRepository, userRepository, passwordService, authService)
	account.NewRoutes(accountService, authService, ginHandler)
	router := ginRouter.Initialize(routerC, ginHandler)

	if err := router.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
