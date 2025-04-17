// Package config provides configuration management for the application
package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	// Telegram configuration
	Telegram struct {
		Token  string
		ChatID string
	}
	
	// HTTP server configuration
	Server struct {
		Port int
	}
}

var (
	// AppConfig is the global application configuration
	AppConfig Config
)

// Initialize sets up the configuration system
func Initialize() {
	// Set up Viper
	viper.SetEnvPrefix("NESTAT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("telegram.token", "")
	viper.SetDefault("telegram.chatid", "")

	// Load config into struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		slog.Error("Failed to unmarshal config", "error", err)
	}

	// Validate required configs
	if AppConfig.Telegram.Token == "" {
		slog.Warn("NESTAT_TELEGRAM_TOKEN environment variable not set - Telegram notifications will be disabled")
	}
	
	if AppConfig.Telegram.ChatID == "" {
		slog.Warn("NESTAT_TELEGRAM_CHATID environment variable not set - Telegram notifications will be disabled")
	}
}

// InitCommands sets up the CLI commands
func InitCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ne-stat-toboy",
		Short: "Web application for the short film 'Не Стать Тобой'",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running server on port", AppConfig.Server.Port)
			// Server is started in main.go
		},
	}

	// Add flags
	rootCmd.PersistentFlags().IntVarP(&AppConfig.Server.Port, "port", "p", 8080, "Port to listen on")
	
	// Bind flags to viper
	if err := viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port")); err != nil {
		slog.Error("Failed to bind flag", "error", err)
	}

	return rootCmd
}