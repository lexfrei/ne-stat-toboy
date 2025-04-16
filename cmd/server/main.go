// Package main provides the entry point for the "Ne Stat Toboy" film website.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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
	h := handler.New()

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Add middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())
	// Enable response compression
	e.Use(echoMiddleware.Gzip())
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
		slog.Info("Server starting", "port", "8080")
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
