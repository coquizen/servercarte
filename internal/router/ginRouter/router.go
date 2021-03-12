package ginRouter

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/logger"
)

// NewGinEngineHandler instantiates a new Gin engine
func NewGinEngineHandler() *gin.Engine {
	// With logger and recovery middlewares
	r := gin.Default()
	return r
}

// Start starts a new http server with Gin act as handler
func Initialize(cfg config.Router, ginEngine *gin.Engine) *http.Server {
	logger.Info.Printf("server starting up at %s:%s", cfg.Host, cfg.Port)
	return &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      ginEngine,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
}
