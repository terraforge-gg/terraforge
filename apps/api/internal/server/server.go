package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/terraforge-gg/terraforge/internal/config"
)

func NewServer(cfg *config.Config, logger *slog.Logger, db *sql.DB) *echo.Echo {
	e := echo.New()

	_ = validator.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendUrl},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"name": "terraforge",
			"env":  cfg.Env,
		})
	})

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":    "ok",
			"timestamp": time.Now().UTC().String(),
		})
	})

	return e
}
