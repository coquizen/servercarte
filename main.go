package main

import (
	"flag"
	"log"

	"github.com/CaninoDev/gastro/server/internal/api/account"
	"github.com/CaninoDev/gastro/server/internal/api/menu"
	"github.com/CaninoDev/gastro/server/internal/api/security"
	"github.com/CaninoDev/gastro/server/internal/api/user"
	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/db/gormDB"
	"github.com/CaninoDev/gastro/server/internal/logger"
	"github.com/CaninoDev/gastro/server/internal/router/ginRouter"
)

var (
	configYAML       = flag.String("c", "config.yml", "configure db")
	populateDatabase = flag.Bool("p", false, "reset and populate the db with sample data")
)

func main() {
	flag.Parse()
	routerC, databaseC, err := config.Load(*configYAML)
	if err != nil {
		logger.Error.Fatalf("error parsing config.yml %v", err)
	}

	gormDB, err := gormDB.Start(databaseC, *populateDatabase)
	if err != nil {
		log.Panic(err)
	}

	ginHandler := ginRouter.NewGinEngineHandler()

	menuRepository := menu.NewGormDBRepository(gormDB)
	menuService := menu.Initialize(menuRepository)
	menu.NewGinRoutes(menuService, ginHandler)

	userRepository := user.NewGormDBRepository(gormDB)
	userService := user.Initialize(userRepository)
	user.NewGinRoutes(userService, ginHandler)

	router := ginRouter.Initialize(routerC, ginHandler)
	router.ListenAndServe()
}
