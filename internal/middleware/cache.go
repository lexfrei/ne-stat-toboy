package middleware

import (
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// CacheControlMiddleware adds cache control headers for Cloudflare
func CacheControlMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			
			// Set cache headers based on content type
			switch {
			case strings.HasPrefix(path, "/static/"):
				ext := filepath.Ext(path)
				switch ext {
				case ".css", ".js", ".jpg", ".jpeg", ".png", ".gif", ".ico", ".svg":
					// Long cache for static assets (30 days)
					c.Response().Header().Set("Cache-Control", "public, max-age=2592000, s-maxage=2592000, stale-while-revalidate=86400")
					// Cloudflare specific cache headers
					c.Response().Header().Set("CDN-Cache-Control", "max-age=2592000")
					c.Response().Header().Set("Cloudflare-CDN-Cache-Control", "max-age=2592000")
				}
			case path == "/" || strings.HasPrefix(path, "/about") || strings.HasPrefix(path, "/team") || 
				strings.HasPrefix(path, "/locations") || strings.HasPrefix(path, "/contact"):
				// Short cache for HTML pages (5 minutes)
				c.Response().Header().Set("Cache-Control", "public, max-age=300, s-maxage=300, stale-while-revalidate=900")
				c.Response().Header().Set("CDN-Cache-Control", "max-age=300")
				c.Response().Header().Set("Cloudflare-CDN-Cache-Control", "max-age=300")
			case strings.HasPrefix(path, "/api/"):
				// No cache for API calls
				c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
				c.Response().Header().Set("Pragma", "no-cache")
				c.Response().Header().Set("Expires", "0")
			}
			
			return next(c)
		}
	}
}