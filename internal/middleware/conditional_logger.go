package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

// ConditionalLogger creates a middleware that logs requests unless they match paths that should be skipped
func ConditionalLogger() echo.MiddlewareFunc {
	config := middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","level":"INFO","msg":"HTTP Request","remote_ip":"${remote_ip}","host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}","status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}","bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02T15:04:05.999999999Z07:00",
		Skipper: func(c echo.Context) bool {
			// Skip logging for health check and metrics endpoints
			path := c.Request().URL.Path
			return strings.HasPrefix(path, "/healthz") || strings.HasPrefix(path, "/metrics")
		},
	}

	return middleware.LoggerWithConfig(config)
}
