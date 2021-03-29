package gormDB

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/CaninoDev/gastro/server/internal/config"
)

// Start returns a configured instance of database{}
func Start(cfg config.Database, populateDatabase bool) (*gorm.DB, error) {
	var gormCfg = &gorm.Config{
		Logger: newLogger(),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Now some conditional settings depending on command line flag
	if populateDatabase {
		gormCfg.DisableForeignKeyConstraintWhenMigrating = true
		db, err := open(cfg, gormCfg)
		if err != nil {
			return db, err
		}
		dbCloser, err := db.DB()
		if err != nil {
			return db, err
		}
		if err := seedDB(db); err != nil {
			return db, err
		}
		if err := dbCloser.Close(); err != nil {
			return db, err
		}
		// Once seeded call function again but with bool set as false
		return Start(cfg, false)
	}

	gormCfg.FullSaveAssociations = true
	gormCfg.AllowGlobalUpdate = true

	return open(cfg, gormCfg)
}

func open(cfg config.Database, gormCfg *gorm.Config) (*gorm.DB, error) {
	var dialect gorm.Dialector

	switch strings.ToLower(cfg.Type) {
	case "mysql":
		dialect = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&allowOldPasswords=1",
			cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name))
	case "postgres":
		dialect = postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name))
	case "sqlite3":
		dialect = sqlite.Open(fmt.Sprintf("file:%s?_auth&_auth_user=%s&_auth_pass=%s",
			cfg.Host, cfg.User, cfg.Pass))
	default:
		return &gorm.DB{}, fmt.Errorf("%s is unsupported", cfg.Type)
	}

	return gorm.Open(dialect, gormCfg)
}

func newLogger() logger.Interface {
	return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 100 * time.Millisecond,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
}
