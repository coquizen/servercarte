package api

import (
	"github.com/CaninoDev/gastro/server/api/menu"
	menuTransport "github.com/CaninoDev/gastro/server/api/menu/transport"
	"github.com/CaninoDev/gastro/server/api/user"
	userTransport "github.com/CaninoDev/gastro/server/api/user/transport"
	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/db"
	"github.com/CaninoDev/gastro/server/internal/router"
)

func Start(cfg *config.RouterConf) error {
	db, err := db.New(cfg)
	if err != nil {
		return err
	}

	srv := router.New()
	menuTransport.NewHTTP(menu.Initialize(db), srv)
	userTransport.NewHTTP(user.Initialize(db), srv)

	if err:= router.Start(srv, cfg); err != nil {
		return err
	}
	return nil
}
