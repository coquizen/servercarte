package main

import (
	"flag"
	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/server"
	"log"
)

var (
	configYAML   = flag.String("c", "config.yml", "configure db")
	seedDatabase = flag.Bool("s", false, "reset and seed the db with sample data")
)

func main() {
	flag.Parse()
	routerC, databaseC, securityC, authC, err := config.Load(*configYAML)
	if err != nil {
		log.Fatalf("error parsing config.yml: %v", err)
	}

	app := server.NewApp(routerC, databaseC, authC, securityC, *seedDatabase)
	app.Run()
}