package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
)

//go:embed .claude/*
var templateFS embed.FS

func main() {
	// Parse command-line flags
	config := parseFlags()

	// Handle version flag
	if config.ShowVersion {
		fmt.Printf("cc-init version %s\n", version)
		os.Exit(0)
	}

	// Handle help flag
	if config.ShowHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create and run the engine
	engine := NewEngine(templateFS, config)
	if err := engine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}