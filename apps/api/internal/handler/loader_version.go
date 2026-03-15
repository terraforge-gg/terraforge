package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/dto"
	custom_errors "github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/service"

	"github.com/labstack/echo/v5"
)

type LoaderVersionHandler struct {
	cfg                  *config.Config
	loaderVersionService service.LoaderVersionService
	logger               *slog.Logger
}

func NewLoaderVersionHandler(cfg *config.Config, logger *slog.Logger, loaderVersionService service.LoaderVersionService) *LoaderVersionHandler {
	return &LoaderVersionHandler{
		cfg:                  cfg,
		loaderVersionService: loaderVersionService,
		logger:               logger,
	}
}

func (h *LoaderVersionHandler) GetLoaderVersionById(c *echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	LoaderVersion, err := h.loaderVersionService.GetLoaderVersionById(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrLoaderVersionNotFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "Loader Version not found.",
			})
		default:
			h.logger.Error("Unhandled get loader version error", "Loader Version ID:", id, "Error:", err)
			return c.JSON(http.StatusInternalServerError, dto.ProblemDetails{
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
		}
	}

	LoaderVersionDto := dto.MapToLoaderVersionResponse(*LoaderVersion)

	return c.JSON(http.StatusOK, LoaderVersionDto)
}

func (h *LoaderVersionHandler) GetLoaderVersions(c *echo.Context) error {
	ctx := c.Request().Context()
	loaderVersions, err := h.loaderVersionService.GetLoaderVersions(ctx)

	if err != nil {
		switch {
		default:
			h.logger.Error("Unhandled get loader versions error", "Error:", err)
			return c.JSON(http.StatusInternalServerError, dto.ProblemDetails{
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
		}
	}

	loaderVersionsDto := make([]dto.LoaderVersionResponse, 0, len(loaderVersions))

	for _, lv := range loaderVersions {
		loaderVersionsDto = append(loaderVersionsDto, dto.MapToLoaderVersionResponse(lv))
	}

	return c.JSON(http.StatusOK, loaderVersionsDto)
}
