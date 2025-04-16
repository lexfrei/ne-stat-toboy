// Package middleware provides custom middleware for the application.
package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

type minifyResponseWriter struct {
	io.Writer
	http.ResponseWriter
	minifier *minify.M
	mediaType string
}

func (mrw *minifyResponseWriter) Write(b []byte) (int, error) {
	if mrw.minifier != nil && mrw.mediaType != "" {
		mb := &bytes.Buffer{}
		if err := mrw.minifier.Minify(mrw.mediaType, mb, bytes.NewReader(b)); err != nil {
			slog.Error("Failed to minify response", "error", err, "contentType", mrw.mediaType)
			return mrw.ResponseWriter.Write(b)
		}
		
		minified := mb.Bytes()
		slog.Debug("Minified response", 
			"contentType", mrw.mediaType, 
			"originalSize", len(b), 
			"minifiedSize", len(minified),
			"reduction", fmt.Sprintf("%.2f%%", (1-float64(len(minified))/float64(len(b)))*100),
		)
		
		// Update content length
		mrw.ResponseWriter.Header().Set(echo.HeaderContentLength, strconv.Itoa(len(minified)))
		return mrw.ResponseWriter.Write(minified)
	}
	return mrw.ResponseWriter.Write(b)
}

// MinifyMiddleware creates a middleware that minifies response bodies by content type.
func MinifyMiddleware() echo.MiddlewareFunc {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^image/svg\\+xml$"), svg.Minify)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			response := c.Response()
			
			// Skip minification for non-GET/POST methods or if header indicates skipping
			if request.Method != http.MethodGet && request.Method != http.MethodPost {
				return next(c)
			}
			
			// Wait for headers to be written
			resWriterBefore := response.Writer
			
			// Process the request
			err := next(c)
			if err != nil {
				return err
			}
			
			// If the response writer has been changed, it means someone else has wrapped it
			if response.Writer != resWriterBefore {
				return nil
			}
			
			// Get content type after handler has set headers
			contentType := response.Header().Get(echo.HeaderContentType)
			if contentType == "" {
				return nil
			}
			
			// Extract the base MIME type
			mediaType, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				slog.Error("Failed to parse content type", "error", err, "contentType", contentType)
				return nil
			}
			
			// Skip binary data, already compressed, or unsupported types
			switch {
			case strings.HasPrefix(mediaType, "image/") && mediaType != "image/svg+xml":
				return nil
			case strings.HasPrefix(mediaType, "audio/"):
				return nil
			case strings.HasPrefix(mediaType, "video/"):
				return nil
			case mediaType == "application/octet-stream":
				return nil
			case mediaType == "application/zip":
				return nil
			case mediaType == "application/x-gzip":
				return nil
			case mediaType == "application/pdf":
				return nil
			}
			
			// Check if minifier exists for this media type
			if _, params, minifierFunc := m.Match(mediaType); minifierFunc == nil {
				slog.Debug("No minifier for content type", "mediaType", mediaType, "params", params)
				return nil
			}
			
			// Wrap the response writer
			response.Writer = &minifyResponseWriter{
				ResponseWriter: response.Writer,
				minifier: m,
				mediaType: mediaType,
			}
			
			return nil
		}
	}
}