package handler

import (
	"regexp"
	"strings"
)

// validateEmail checks if the email format is valid
func validateEmail(email string) bool {
	// Simple validation just for demonstration - using raw string
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// sanitizeString removes potentially harmful characters
func sanitizeString(s string) string {
	// Trim whitespace and normalize
	s = strings.TrimSpace(s)

	// Replace multiple spaces with a single space - using raw string
	spaceRegex := regexp.MustCompile(`\s+`)
	s = spaceRegex.ReplaceAllString(s, " ")

	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")

	return s
}
