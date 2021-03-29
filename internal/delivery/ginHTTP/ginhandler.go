package ginHTTP

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/logger"
)

// NewServer starts a new ginHTTP server with Gin acting as the handler for routes
func NewServer(cfg config.Router, ginHandler *gin.Engine) *http.Server {
	logger.Info.Printf("server starting up at %s:%s", cfg.Host, cfg.Port)
	return &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      ginHandler,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
}

// NewHandler instantiates a new Gin engine
func NewHandler(rCfg config.Router) *gin.Engine {
	// With logger and recovery middlewares
	r := gin.Default()
	return r
}