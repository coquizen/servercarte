package database

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/CaninoDev/gastro/server/config"
	"github.com/CaninoDev/gastro/server/models"
)

//New returns a configured instance of database{}
func New(cfg *config.Config) (*gorm.DB, error) {

	var dialect gorm.Dialector

	switch strings.ToLower(cfg.Database.Type) {
	case "mysql":
		dialect = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&allowOldPasswords=1",
			cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name))
	case "postgres":
		dialect = postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name))
	case "sqlite3":
		dialect = sqlite.Open(fmt.Sprintf("file:%s?_auth&_auth_user=%s&_auth_pass=%s",
			cfg.Database.Host, cfg.Database.User, cfg.Database.Pass))
	default:
		return &gorm.DB{}, fmt.Errorf("%s is unsupported", cfg.Database.Type)
	}
		gormDB, err := gorm.Open(dialect, config.Gorm)
		if err != nil {
			return gormDB, err
		}

		// Update the db's schema with any changes
		if err := gormDB.AutoMigrate(models.Section{}, models.Item{}, models.User{}); err != nil {
			return gormDB, err
		}
		return gormDB, nil
}
