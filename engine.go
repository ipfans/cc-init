package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// Engine is the core orchestrator for the cc-init tool
type Engine struct {
	templateFS embed.FS
	config     *Config
	logger     *Logger
	fs         FileSystem
	tmpl       *TemplateManager
	stats      Statistics
}

// Statistics tracks the operation results
type Statistics struct {
	FilesCreated      int
	FilesSkipped      int
	DirsCreated       int
	DirsSkipped       int
	Errors            []error
}

// NewEngine creates a new Engine instance
func NewEngine(templateFS embed.FS, config *Config) *Engine {
	logger := NewLogger(config.Verbose, config.NoColor)
	
	var fileSystem FileSystem = NewOSFileSystem()
	if config.DryRun {
		fileSystem = NewDryRunFileSystem(fileSystem, logger)
	}
	
	return &Engine{
		templateFS: templateFS,
		config:     config,
		logger:     logger,
		fs:         fileSystem,
		tmpl:       NewTemplateManager(templateFS, ".claude"),
		stats:      Statistics{},
	}
}

// Run executes the main initialization process
func (e *Engine) Run() error {
	e.logger.Debug("Starting cc-init with target directory: %s", e.config.TargetDir)
	
	// Check if templates exist
	if !e.tmpl.HasTemplates() {
		return fmt.Errorf("no template files found in embedded .claude directory")
	}
	
	// List all templates if verbose
	if e.config.Verbose {
		templates, err := e.tmpl.ListTemplates()
		if err == nil {
			e.logger.Debug("Found %d template files", len(templates))
			for _, tmpl := range templates {
				e.logger.Debug("  - %s", tmpl)
			}
		}
	}
	
	// Walk through all template files
	err := e.tmpl.Walk(func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			e.logger.Error("Error accessing %s: %v", path, err)
			e.stats.Errors = append(e.stats.Errors, err)
			return nil // Continue processing other files
		}
		
		targetPath := filepath.Join(e.config.TargetDir, ".claude", filepath.FromSlash(path))
		
		if entry.IsDir() {
			return e.processDirectory(targetPath)
		}
		
		return e.processFile(path, targetPath)
	})
	
	if err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}
	
	// Show summary
	e.showSummary()
	
	// Return error if there were any critical errors
	if len(e.stats.Errors) > 0 {
		return fmt.Errorf("completed with %d errors", len(e.stats.Errors))
	}
	
	return nil
}

// processDirectory handles directory creation
func (e *Engine) processDirectory(targetPath string) error {
	e.logger.Debug("Processing directory: %s", targetPath)
	
	if e.fs.Exists(targetPath) {
		info, err := e.fs.Stat(targetPath)
		if err != nil {
			return fmt.Errorf("failed to stat %s: %w", targetPath, err)
		}
		
		if !info.IsDir() {
			err := fmt.Errorf("path exists but is not a directory: %s", targetPath)
			e.stats.Errors = append(e.stats.Errors, err)
			return err
		}
		
		e.logger.DirSkipped(e.formatPath(targetPath))
		e.stats.DirsSkipped++
		return nil
	}
	
	err := e.fs.CreateDir(targetPath, e.tmpl.GetDefaultDirMode())
	if err != nil {
		e.logger.Error("Failed to create directory %s: %v", targetPath, err)
		e.stats.Errors = append(e.stats.Errors, err)
		return err
	}
	
	e.logger.DirCreated(e.formatPath(targetPath))
	e.stats.DirsCreated++
	return nil
}

// processFile handles file creation
func (e *Engine) processFile(sourcePath, targetPath string) error {
	e.logger.Debug("Processing file: %s -> %s", sourcePath, targetPath)
	
	if e.fs.Exists(targetPath) {
		info, err := e.fs.Stat(targetPath)
		if err != nil {
			return fmt.Errorf("failed to stat %s: %w", targetPath, err)
		}
		
		if info.IsDir() {
			err := fmt.Errorf("path exists but is a directory: %s", targetPath)
			e.stats.Errors = append(e.stats.Errors, err)
			return err
		}
		
		e.logger.FileSkipped(e.formatPath(targetPath))
		e.stats.FilesSkipped++
		return nil
	}
	
	// Read the source file
	content, err := e.tmpl.ReadFile(sourcePath)
	if err != nil {
		e.logger.Error("Failed to read template file %s: %v", sourcePath, err)
		e.stats.Errors = append(e.stats.Errors, err)
		return err
	}
	
	// Get file mode
	mode := e.tmpl.GetDefaultFileMode(sourcePath)
	
	// Create the file
	err = e.fs.CreateFile(targetPath, content, mode)
	if err != nil {
		e.logger.Error("Failed to create file %s: %v", targetPath, err)
		e.stats.Errors = append(e.stats.Errors, err)
		return err
	}
	
	e.logger.FileCreated(e.formatPath(targetPath))
	e.stats.FilesCreated++
	return nil
}

// formatPath formats a path for display
func (e *Engine) formatPath(path string) string {
	// Try to make path relative to target directory for cleaner output
	if relPath, err := filepath.Rel(e.config.TargetDir, path); err == nil && !strings.HasPrefix(relPath, "..") {
		return relPath
	}
	return path
}

// showSummary displays the operation summary
func (e *Engine) showSummary() {
	totalCreated := e.stats.FilesCreated + e.stats.DirsCreated
	totalSkipped := e.stats.FilesSkipped + e.stats.DirsSkipped
	
	fmt.Println() // Empty line before summary
	
	if e.config.DryRun {
		e.logger.Info("DRY RUN - No changes were made")
		fmt.Println()
	}
	
	// Show what was created
	if totalCreated > 0 {
		items := []string{}
		if e.stats.FilesCreated > 0 {
			items = append(items, fmt.Sprintf("%d %s", e.stats.FilesCreated, pluralize("file", e.stats.FilesCreated)))
		}
		if e.stats.DirsCreated > 0 {
			items = append(items, fmt.Sprintf("%d %s", e.stats.DirsCreated, pluralize("directory", e.stats.DirsCreated)))
		}
		e.logger.Success("Created %s", strings.Join(items, " and "))
	}
	
	// Show what was skipped
	if totalSkipped > 0 {
		items := []string{}
		if e.stats.FilesSkipped > 0 {
			items = append(items, fmt.Sprintf("%d %s", e.stats.FilesSkipped, pluralize("file", e.stats.FilesSkipped)))
		}
		if e.stats.DirsSkipped > 0 {
			items = append(items, fmt.Sprintf("%d %s", e.stats.DirsSkipped, pluralize("directory", e.stats.DirsSkipped)))
		}
		e.logger.Info("Skipped %s (already exist)", strings.Join(items, " and "))
	}
	
	// Show errors if any
	if len(e.stats.Errors) > 0 {
		e.logger.Error("Encountered %d %s during initialization", len(e.stats.Errors), pluralize("error", len(e.stats.Errors)))
		if e.config.Verbose {
			for _, err := range e.stats.Errors {
				e.logger.Error("  - %v", err)
			}
		}
	}
	
	// Final status
	if totalCreated == 0 && totalSkipped > 0 {
		e.logger.Info("All Claude configuration files already exist")
	} else if len(e.stats.Errors) == 0 {
		e.logger.Success("Claude configuration initialized successfully")
	}
}

// pluralize returns the plural form of a word if count != 1
func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	
	// Handle special cases
	switch word {
	case "directory":
		return "directories"
	default:
		return word + "s"
	}
}