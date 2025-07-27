package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const version = "0.1.0"

// Config holds the CLI configuration
type Config struct {
	TargetDir   string
	DryRun      bool
	Verbose     bool
	NoColor     bool
	ShowHelp    bool
	ShowVersion bool
}

// parseFlags parses command-line flags and returns the configuration
func parseFlags() *Config {
	config := &Config{}

	// Define flags
	flag.StringVar(&config.TargetDir, "target", ".", "Target directory for initialization")
	flag.StringVar(&config.TargetDir, "t", ".", "Target directory for initialization (shorthand)")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Preview operations without making changes")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&config.Verbose, "v", false, "Enable verbose output (shorthand)")
	flag.BoolVar(&config.NoColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version information")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "cc-init - Initialize Claude Code configuration\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # Initialize in current directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t ./myproject     # Initialize in ./myproject\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --dry-run          # Preview what would be created\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -v                 # Show detailed output\n", os.Args[0])
	}

	// Parse flags
	flag.Parse()

	// Handle help flag
	if flag.NArg() > 0 && flag.Arg(0) == "help" {
		config.ShowHelp = true
	}

	return config
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Convert target directory to absolute path
	absPath, err := filepath.Abs(config.TargetDir)
	if err != nil {
		return fmt.Errorf("invalid target directory: %w", err)
	}
	config.TargetDir = absPath

	// Check if target directory exists
	info, err := os.Stat(config.TargetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("target directory does not exist: %s", config.TargetDir)
		}
		return fmt.Errorf("failed to access target directory: %w", err)
	}

	// Ensure it's a directory
	if !info.IsDir() {
		return fmt.Errorf("target path is not a directory: %s", config.TargetDir)
	}

	// Check write permissions (unless dry-run)
	if !config.DryRun {
		// Try to create a temporary file to test write permissions
		testFile := filepath.Join(config.TargetDir, ".cc-init-test")
		f, err := os.Create(testFile)
		if err != nil {
			return fmt.Errorf("no write permission in target directory: %w", err)
		}
		f.Close()
		os.Remove(testFile)
	}

	return nil
}