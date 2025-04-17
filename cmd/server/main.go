// Package main provides the entry point for the "Ne Stat Toboy" film website.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/lexfrei/ne-stat-toboy/internal/config"
	"github.com/lexfrei/ne-stat-toboy/internal/handler"
	"github.com/lexfrei/ne-stat-toboy/internal/middleware"
	"github.com/lexfrei/ne-stat-toboy/internal/minify"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize configuration
	config.Initialize()
	rootCmd := config.InitCommands()
	
	// If called with arguments, let cobra handle it
	if len(os.Args) > 1 {
		if err := rootCmd.Execute(); err != nil {
			slog.Error("Command execution failed", "error", err)
			os.Exit(1)
		}
		return
	}

	// Create main context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Minify static files
	staticDir := "web/static"
	slog.Info("Minifying static files", "directory", staticDir)
	if err := minify.MinifyStaticFiles(staticDir); err != nil {
		slog.Error("Failed to minify static files", "error", err)
		// Continue execution even if minification fails
	}

	// Setup handlers
	h := handler.New(
		handler.WithTelegramConfig(
			config.AppConfig.Telegram.Token,
			config.AppConfig.Telegram.ChatID,
		),
	)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Add middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())
	// Security middlewares
	e.Use(echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		TokenLookup: "form:_csrf",
		CookieName: "csrf",
		CookieMaxAge: 3600,
		CookieSecure: true, 
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
		ContextKey: "csrf",
	}))
	e.Use(echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' unpkg.com; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self' api.telegram.org; font-src 'self'; base-uri 'self'; form-action 'self'; frame-ancestors 'self'",
	}))
	// Rate limiting
	e.Use(echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
		Store: echoMiddleware.NewRateLimiterMemoryStore(20), // 20 requests per second
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "too many requests",
			})
		},
	}))
	// Enable response compression
	e.Use(echoMiddleware.Gzip())
	// Cache control for Cloudflare
	e.Use(middleware.CacheControlMiddleware())
	// Add minification middleware
	e.Use(middleware.MinifyMiddleware())

	// Static files handler
	e.Static("/static", staticDir)

	// Setup routes
	e.GET("/", h.HomeHandlerEcho)
	e.GET("/about", h.AboutHandlerEcho)
	e.GET("/team", h.TeamHandlerEcho)
	e.GET("/locations", h.LocationsHandlerEcho)
	e.GET("/contact", h.ContactHandlerEcho)
	e.POST("/api/contact", h.ContactSubmitHandlerEcho)

	// Start server in a goroutine
	go func() {
		port := config.AppConfig.Server.Port
		address := fmt.Sprintf(":%d", port)
		slog.Info("Server starting", "port", port)
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Create a deadline to wait for
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	slog.Info("Server stopped")
}