package main

import (
	"log"
	"os"

	"github.com/lexfrei/ne-stat-toboy/internal/minify"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <static-dir>", os.Args[0])
	}

	staticDir := os.Args[1]
	if err := minify.MinifyStaticFiles(staticDir); err != nil {
		log.Fatalf("Error minifying files: %v", err)
	}
}
