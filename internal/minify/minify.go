// Package minify provides utilities for minifying static files.
package minify

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

// MinifyStaticFiles minifies CSS, JS, HTML and SVG files in the specified directory.
func MinifyStaticFiles(staticDir string) error {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	// Create a standalone command for minification
	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		var mediaType string
		switch ext {
		case ".css":
			mediaType = "text/css"
		case ".js":
			mediaType = "application/javascript"
		case ".html", ".htm":
			mediaType = "text/html"
		case ".svg":
			mediaType = "image/svg+xml"
		default:
			return nil
		}

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		if len(content) == 0 {
			return nil
		}

		// Check if minifier exists for this type
		_, _, minifierFunc := m.Match(mediaType)
		if minifierFunc == nil {
			slog.Debug("No minifier for file type", "path", path, "mediaType", mediaType)
			return nil
		}

		// Minify content
		minified, err := m.Bytes(mediaType, content)
		if err != nil {
			return fmt.Errorf("failed to minify %s: %w", path, err)
		}

		if len(minified) >= len(content) {
			slog.Debug("No size reduction after minification", "file", path)
			return nil
		}

		// Write minified content
		if err := os.WriteFile(path, minified, info.Mode()); err != nil {
			return fmt.Errorf("failed to write minified file %s: %w", path, err)
		}

		slog.Info("Minified file", 
			"file", path,
			"originalSize", len(content),
			"minifiedSize", len(minified),
			"reduction", fmt.Sprintf("%.2f%%", (1-float64(len(minified))/float64(len(content)))*100),
		)

		return nil
	})
}