package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/logger"
	"github.com/terraforge-gg/terraforge/internal/server"
)

func main() {
	logger := logger.New()
	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg.DatabaseUrl)

	if err != nil {
		logger.Error("An error occurred connecting to the database", "error", err)
		return
	}

	defer db.Close()

	err = database.Migrate(logger, db)

	if err != nil {
		logger.Error("An error occurred during migrate", "error", err)
		return
	}

	e, err := server.NewServer(cfg, logger, db)

	if err != nil {
		logger.Error("Failed to create server", "error", err)
		return
	}

	go func() {
		err := e.Start(":" + cfg.HostPort)

		if err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("Server gracefully stopped")
}
