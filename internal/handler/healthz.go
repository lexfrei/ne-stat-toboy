package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// HealthCheckHandler handles health check requests
func (h *Handler) HealthCheckHandler(c echo.Context) error {
	// Simple health check - just return OK
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
